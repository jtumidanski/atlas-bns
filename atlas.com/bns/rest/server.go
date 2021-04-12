package rest

import (
	"atlas-bns/name"
	"log"
	"net/http"
	"os"
	"time"
)
import "github.com/gorilla/mux"

type server struct {
	l  *log.Logger
	hs *http.Server
}

func GetServer(l *log.Logger) *server {
	r := mux.NewRouter().PathPrefix("/ms/bns").Subrouter()
	r.Use(commonHeader)

	cr := r.PathPrefix("/names").Subrouter()
	cr.HandleFunc("", name.GetName(l)).Methods(http.MethodGet).Queries("name", "{name}")
	cr.HandleFunc("", name.GetNames(l)).Methods(http.MethodGet)

	hs := http.Server{
		Addr:         ":8080",
		Handler:      r,
		ErrorLog:     l,                 // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}
	return &server{l, &hs}
}

func (s *server) Run() {
	s.l.Println("[INFO] Starting server on port 8080")
	err := s.hs.ListenAndServe()
	if err != nil {
		s.l.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}

func commonHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
