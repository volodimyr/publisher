package response

import (
	"github.com/volodimyr/event_publisher/api/config"
	"net/http"
)

func OK(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set(config.ContentType, config.CharSet)
	w.Write([]byte("OK"))
}

func Created(w http.ResponseWriter) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set(config.ContentType, config.CharSet)
	w.Write([]byte("Resource is created"))
}
