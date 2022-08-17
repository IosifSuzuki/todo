package handler

import "net/http"

type CommonHanlder struct{}

func (c CommonHanlder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadGateway)
}
