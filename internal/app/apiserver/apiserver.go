package apiserver

import (
	"net/http"
	"refactoring/internal/app/store/jsonstore"
)

func Start(config *Config) error {
	store := jsonstore.New(config.UserStoreUrl)
	srv := newServer(store)

	return http.ListenAndServe(config.BindAddr, srv)
}
