package main

import (
	"net/http"
	"strings"
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
		if route.method == req.Method && route.path == req.URL.Path {
			handler = route.handler
		}
	}

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

type routeTree struct {
	root *routeNode
}

func newRouteTree(path string) *routeTree {
	root := &routeNode{
		value:    path,
		handlers: make(map[string]http.HandlerFunc),
	}

	return &routeTree{root}
}

func (rt *routeTree) Get(path string, handler http.HandlerFunc) {
	if len(path) == 1 {
		rt.root = newRouteNode(path)
		return
	}

	rt.root = newRouteNode("/")
	tokens := strings.Split(path, "/")

	current := rt.root

	for i := 1; i < len(tokens); i++ {
		node := newRouteNode(tokens[i])
		current.child = node
		current = node
	}

	current.handlers["GET"] = handler
}

func (rt *routeTree) Post(path string, handler http.HandlerFunc) {
}

type routeNode struct {
	value    string
	handlers map[string]http.HandlerFunc
	child    *routeNode
}

func newRouteNode(value string) *routeNode {
	return &routeNode{
		value:    value,
		handlers: make(map[string]http.HandlerFunc),
	}
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
