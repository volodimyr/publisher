package resp

import "net/http"

//OK uses to stablish OK response
func OK(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-type", "text/plain; charset=utf-8")
	w.Write([]byte(msg))
}

//Created uses to notify client 'resource has been created'
func Created(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-type", "text/plain; charset=utf-8")
	w.Write([]byte(msg))
}
