package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var config *envConfig

func init() {
	config = newEnvConfig()
}

func main() {
	r := mux.NewRouter()

	srv := &http.Server{
		Addr:         "127.0.0.1:" + config.Port,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		Handler:      r,
	}

	redirectTmpl := template.Must(template.ParseFiles("./web/templates/redirect.html"))
	successTmpl := template.Must(template.ParseFiles("./web/templates/success.html"))

	r.HandleFunc("/auth/callback", handleOauthCallback(successTmpl)).Methods("GET")
	r.HandleFunc("/auth", handleClientOauth(redirectTmpl)).Methods("GET")

	fmt.Println("Starting web server...http://127.0.0.1:" + config.Port)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("[Error] --server: ", err.Error())
	}
}

func handleClientOauth(tmpl *template.Template) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		clientID := config.OauthClientID
		scope := "read-protected write-protected"
		state := config.CsrfToken
		uri := fmt.Sprintf("https://id.heroku.com/oauth/authorize?client_id=%s&response_type=code&scope=%s&state=%s", clientID, scope, state)

		data := map[string]interface{}{
			"authUrl": uri,
		}

		w.Header().Set("Content-type", "text/html")
		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func handleOauthCallback(tmpl *template.Template) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query()["code"][0]
		state := r.URL.Query()["state"][0]
		var tokenData map[string]interface{}

		if state != config.CsrfToken {
			http.Error(w, errors.New("authorization failed").Error(), http.StatusUnauthorized)
			return
		}

		// handle token exchangze
		uri := fmt.Sprintf("https://id.heroku.com/oauth/token")

		d := url.Values{}
		d.Set("grant_type", "authorization_code")
		d.Set("code", code)
		d.Set("client_secret", config.OauthSecret)

		client := &http.Client{}

		req, err := http.NewRequest("POST", uri, strings.NewReader(d.Encode()))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		defer resp.Body.Close()

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.Unmarshal(b, &tokenData); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, ok := tokenData["access_token"]; !ok {
			http.Error(w, errors.New("authorization failed").Error(), http.StatusInternalServerError)
			return
		}

		uri = fmt.Sprintf("http://localhost:7070/auth/callback")

		data, err := json.Marshal(map[string]string{
			"access_token":  tokenData["access_token"].(string),
			"refresh_token": tokenData["refresh_token"].(string),
		})

		body := bytes.NewBuffer(data)

		req, err = http.NewRequest("POST", uri, body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, err := client.Do(req); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-type", "text/html")
		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

type envConfig struct {
	Port          string `mapstructure:"PORT"`
	OauthClientID string `mapstructure:"OAUTH_CLIENT_ID"`
	OauthSecret   string `mapstructure:"OAUTH_SECRET"`
	CsrfToken     string `mapstructure:"CSRF_TOKEN"`
}

func newEnvConfig() *envConfig {
	var config envConfig

	if err := godotenv.Load(); err != nil {
		log.Fatalf("[Error] --config: couldn't load env file - %s", err.Error())
	}

	config.Port = os.Getenv("PORT")
	config.OauthClientID = os.Getenv("OAUTH_CLIENT_ID")
	config.OauthSecret = os.Getenv("OAUTH_SECRET")
	config.CsrfToken = os.Getenv("CSRF_TOKEN")

	return &config
}
