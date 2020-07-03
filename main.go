package main

import "net/http"

type handler struct{}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hi from the outer space"))
}

func main() {
	http.ListenAndServe(":9090", handler{})
}
