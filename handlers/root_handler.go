package handlers

import (
	"net/http"
)

// RootHandler handles the "/" route
func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to the Book API!\nTry visiting '/books' to get started."))
}
