package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/Mayowa-Ojo/otter/internal"
)

// AuthorizeClient - grant client access via heroku oauth
func AuthorizeClient() error {
	url := fmt.Sprintf("http://localhost:5000/auth")
	timeout := time.NewTimer(30 * time.Second)

	if err := internal.OpenURLLink(url); err != nil {
		return err
	}

	shutdownChan := make(chan string)

	callbackHandler := func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err.Error())
		}

		if err = json.Unmarshal(b, &body); err != nil {
			log.Fatal(err.Error())
		}

		if _, ok := body["access_token"]; !ok {
			log.Fatal(errors.New("missing credentials"))
		}

		if err := internal.PersistAuthorization(body["access_token"].(string), body["refresh_token"].(string)); err != nil {
			log.Fatal(err.Error())
		}

		fmt.Println(w, "\nYou are now logged in \u2713")
		shutdownChan <- "done"
		timeout.Stop()
	}

	handlers := []internal.HTTPHandler{
		{
			Handler: callbackHandler,
			Path:    "/auth/callback",
			Method:  "POST",
		},
	}

	// kill the server after 30s if no response is recieved
	go func(t *time.Timer, ch chan string) {
		<-t.C
		ch <- "done"
	}(timeout, shutdownChan)

	internal.EphemeralServer(shutdownChan, handlers)

	return nil
}
