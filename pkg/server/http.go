package server

import (
	"net/http"
	"time"
)

//New makes custom server configuration
func New(sm *http.ServeMux, addr string) *http.Server {
	return &http.Server{
		Addr:         addr,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 10,
		IdleTimeout:  time.Second * 120,
		Handler:      sm,
	}
}
