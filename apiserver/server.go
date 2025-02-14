package apiserver

import (
	"context"
	"net/http"

	"github.com/Jidetireni/async-api.git/config"
)

type ApiServer struct {
	Conf *config.Config
}

func New(conf *config.Config) *ApiServer {
	return &ApiServer{
		Conf: conf,
	}
}

func (s *ApiServer) ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))

}

func (s *ApiServer) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", s.ping)
	router := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	return router.ListenAndServe()
}
