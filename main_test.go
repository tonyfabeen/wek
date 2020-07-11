package main

import (
	"bytes"
	"net/http"
	"testing"
)

var defaultHandler = func(message string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(message))
	}
}

type mockResponseWriter struct {
	writerBuffer *bytes.Buffer
}

func newMockResponseWriter() *mockResponseWriter {
	return &mockResponseWriter{
		writerBuffer: bytes.NewBuffer([]byte("")),
	}
}

func (mrw mockResponseWriter) Header() http.Header {
	return nil
}
func (mrw mockResponseWriter) Write(content []byte) (int, error) {
	len, _ := mrw.writerBuffer.Write(content)

	return len, nil
}
func (mrw mockResponseWriter) WriteHeader(statusCode int) {}

func TestAddRoute(t *testing.T) {
	t.Run("when no root", func(t *testing.T) {
		t.Run("and add root", func(t *testing.T) {
			routeTree := routeTree{}
			routeTree.Add(http.MethodGet, "/", defaultHandler("GET / handler"))

			if routeTree.root == nil {
				t.Error("should have a root node")
			}

			if len(routeTree.root.handlers) == 0 {
				t.Error("root node should have handlers")
			}

			if handler := routeTree.root.handlers[http.MethodGet]; handler == nil {
				t.Error("root node should have the handler")
			}
		})

		t.Run("and add root and child node", func(t *testing.T) {
			routeTree := routeTree{}
			routeTree.Add(http.MethodGet, "/posts", defaultHandler("GET /posts handler"))

			if routeTree.root == nil {
				t.Error("should have a root node")
			}

			if len(routeTree.root.handlers) != 0 {
				t.Error("root node should not have handlers")
			}

			if routeTree.root.child == nil {
				t.Error("should have a child")
			}

			if routeTree.root.child.value != "posts" {
				t.Error("should have the correct value")
			}

			if handler := routeTree.root.child.handlers[http.MethodGet]; handler == nil {
				t.Error("child node should have the handler")
			}
		})
	})

	t.Run("when root", func(t *testing.T) {
		t.Run("and add a handler to root", func(t *testing.T) {
			routeTree := routeTree{}
			routeTree.Add(http.MethodGet, "/", defaultHandler("GET / handler"))
			routeTree.Add(http.MethodPost, "/", defaultHandler("POST / handler"))

			if len(routeTree.root.handlers) != 2 {
				t.Error("should have both handlers")
			}

			if handler := routeTree.root.handlers["GET"]; handler == nil {
				t.Error("should have a GET handler")
			}

			if handler := routeTree.root.handlers["POST"]; handler == nil {
				t.Error("should have a POST handler")
			}
		})

		t.Run("and add a child node", func(t *testing.T) {
			routeTree := routeTree{}
			routeTree.Add(http.MethodGet, "/", defaultHandler("GET / handler"))
			routeTree.Add(http.MethodGet, "/posts", defaultHandler("GET /posts handler"))

			if len(routeTree.root.handlers) != 1 {
				t.Error("should have a handler")
			}

			if handler := routeTree.root.handlers["GET"]; handler == nil {
				t.Error("should have a GET handler")
			}

			if routeTree.root.child.value != "posts" {
				t.Error("should have the correct value")
			}

			if handler := routeTree.root.child.handlers["GET"]; handler == nil {
				t.Error("child should have a GET handler")
			}
		})
	})
}

func TestFindRoute(t *testing.T) {
	routeTree := routeTree{}
	routeTree.Add(http.MethodGet, "/", defaultHandler("GET / handler"))

	node := routeTree.Find("/")
	if node.value != "/" {
		t.Error("should return the correct node")
	}

	routeTree.Add(http.MethodGet, "/posts", defaultHandler("GET /posts handler"))
	node = routeTree.Find("/posts")
	if node.value != "posts" {
		t.Error("should return the correct node")
	}

	node = routeTree.Find("/authors")
	if node != nil {
		t.Error("should return nil node")
	}

	routeTree.Add(http.MethodGet, "/posts/authors", defaultHandler("GET /authors handler"))
	node = routeTree.Find("/posts/authors")
	if node.value != "authors" {
		t.Error("should return the correct node")
	}
}
