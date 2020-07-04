package main

import (
	"log"
	"net/http"
)

type router struct {
	routes []route
}

func newRouter() *router {
	routes := []route{}
	return &router{routes}
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler := r.match(req)
	handler(w, req)
}

func (r *router) match(req *http.Request) http.HandlerFunc {
	handler := noRoutesMatches

	for _, route := range r.routes {
		if route.method == req.Method {
			handler = route.handler
		}
	}

	log.Printf("handler %T", handler)
	return handler
}

func (r *router) Get(path string, handler http.HandlerFunc) {
	route := route{
		method:  http.MethodGet,
		path:    path,
		handler: handler,
	}

	r.routes = append(r.routes, route)
}

func (r *router) Post(path string, handler http.HandlerFunc) {
	route := route{
		method:  http.MethodPost,
		path:    path,
		handler: handler,
	}

	r.routes = append(r.routes, route)
}

type route struct {
	method  string
	path    string
	handler http.HandlerFunc
}

func noRoutesMatches(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("No route matches"))
}

func main() {
	router := newRouter()
	router.Get("/", func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodGet {
			w.Write([]byte("handle GET method"))
		}
	})
	router.Post("/", func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPost {
			w.Write([]byte("handle POST method"))
		}
	})

	http.ListenAndServe(":9090", router)
}
