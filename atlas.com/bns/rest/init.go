package rest

import (
	"atlas-bns/name"
	"context"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

func CreateRestService(l *logrus.Logger, ctx context.Context, wg *sync.WaitGroup) {
	go NewServer(l, ctx, wg, ProduceRoutes)
}

func ProduceRoutes(l logrus.FieldLogger) http.Handler {
	r := mux.NewRouter().PathPrefix("/ms/bns").Subrouter()
	r.Use(CommonHeader)

	cr := r.PathPrefix("/names").Subrouter()
	cr.HandleFunc("", name.GetName(l)).Methods(http.MethodGet).Queries("name", "{name}")
	cr.HandleFunc("", name.GetNames(l)).Methods(http.MethodGet)

	return r
}
