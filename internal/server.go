package internal

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// HTTPHandler -
type HTTPHandler struct {
	Handler func(http.ResponseWriter, *http.Request)
	Method  string
	Path    string
}

// EphemeralServer - create a short-lived server to listen for requests
func EphemeralServer(shutdownChan chan string, handlers []HTTPHandler) {
	r := mux.NewRouter()

	srv := &http.Server{
		Addr:         "127.0.0.1:7070",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		Handler:      r,
	}

	for _, h := range handlers {
		r.HandleFunc(h.Path, h.Handler).Methods(h.Method)
	}

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("up and running..."))

		shutdownChan <- "done"
	}).Methods("POST")

	go func(s *http.Server) {
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatal("[Error] --web: ", err.Error())
		}
	}(srv)

	msg := <-shutdownChan

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	if msg != "done" {
		log.Fatal("[Error]: something went wrong")
		srv.Shutdown(ctx)
	}
}
