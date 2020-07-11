package main

import (
	"net/http"
	"strings"
)

type router struct {
	*routeTree
}

func newRouter() *router {
	return &router{}
}

func (r *router) Get(path string, handler http.HandlerFunc) {
	r.routeTree.Add(http.MethodGet, path, handler)
}

func (r *router) Post(path string, handler http.HandlerFunc) {
	r.routeTree.Add(http.MethodPost, path, handler)
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// handler := r.match(req)
	// handler(w, req)
}

type routeTree struct {
	root *routeNode
}

func (rt *routeTree) Add(httpMethod string, path string, handler http.HandlerFunc) {
	if rt.root == nil {
		rt.root = newRouteNode("/")
	}

	if path == "/" {
		rt.root.handlers[httpMethod] = handler
		return
	}

	pathTokens := strings.Split(path, "/")
	node := rt.root

	for _, pathToken := range pathTokens[1:] {
		if node.child == nil {
			node.child = newRouteNode(pathToken)
		}
		node = node.child
	}

	node.handlers[httpMethod] = handler
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
		w.Write([]byte("handle GET method"))
	})
	router.Post("/", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("handle POST method"))
	})

	http.ListenAndServe(":9090", router)
}
