package handler

import (
	"net/http"
	"strings"
)

// DefaultHandler ...
type DefaultHandler struct{}

// NewDefaultHandler ...
func NewDefaultHandler() *DefaultHandler {
	h := &DefaultHandler{}
	return h
}

func (h *DefaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte(strings.Join([]string{"Default response to: ", "(" + r.Method + ")", r.URL.Path}, " ")))
}
