package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

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

func TestRouteTree(t *testing.T) {
	tree := newRouteTree("/")

	root := tree.root
	if root == nil {
		t.Error("should have a root node")
	}

	if root.value != "/" {
		t.Error("should have the right value")
	}

	if root.child != nil {
		t.Error("should not have a child")
	}

	tree = &routeTree{}
	tree.Get("/", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("handle GET method"))
	})

	root = tree.root
	if root == nil {
		t.Error("should have a root node")
	}

	if root.value != "/" {
		t.Error("should have the right value")
	}

	child := root.child
	if child != nil {
		t.Errorf("should not have a child %v", child)
	}

	tree = &routeTree{}
	tree.Get("/posts", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("handle GET /posts"))
	})

	root = tree.root
	if root == nil {
		t.Error("should have a root node")
	}

	if root.value != "/" {
		t.Error("should have the right value")
	}

	child = root.child
	if child == nil {
		t.Error("should have a child")
	}

	if child.handlers["GET"] == nil {
		t.Error("should have a GET handler")
	}

	if child.handlers["POST"] != nil {
		t.Error("should not have a POST handler")
	}
}

func TestMatch(t *testing.T) {
	t.Run("when POST method", func(t *testing.T) {
		postHandler := func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte("handle POST method\n"))
		}

		router := newRouter()
		router.Post("/", postHandler)

		writer := newMockResponseWriter()
		requestURL := &url.URL{Path: "/"}
		request := http.Request{Method: http.MethodPost, URL: requestURL}
		handler := router.match(&request)

		postHandler(writer, &request)
		handler(writer, &request)

		expected := "handle POST method\nhandle POST method\n"
		got := writer.writerBuffer.String()
		if expected != got {
			t.Errorf("should return the right content\nexpected: %s \ngot: %s", expected, got)
		}
	})

	t.Run("when GET method", func(t *testing.T) {
		getHandler := func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte("handle GET method\n"))
		}

		router := newRouter()
		router.Get("/", getHandler)

		writer := newMockResponseWriter()
		requestURL := &url.URL{Path: "/"}
		request := http.Request{Method: http.MethodGet, URL: requestURL}
		handler := router.match(&request)

		getHandler(writer, &request)
		handler(writer, &request)

		expected := "handle GET method\nhandle GET method\n"
		got := writer.writerBuffer.String()
		if expected != got {
			t.Errorf("should return the right content\nexpected: %s \ngot: %s", expected, got)
		}
	})

	t.Run("when no route matches", func(t *testing.T) {
		router := newRouter()

		writer := newMockResponseWriter()
		request := http.Request{Method: http.MethodGet}
		handler := router.match(&request)

		handler(writer, &request)

		expected := "No route matches"
		got := writer.writerBuffer.String()
		if expected != got {
			t.Errorf("should return the right content\nexpected: %s \ngot: %s", expected, got)
		}
	})
}

func TestRequests(t *testing.T) {
	router := newRouter()
	getHandler := func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("handle GET method"))
	}
	router.Get("/", getHandler)

	postHandler := func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("handle POST method"))
	}
	router.Post("/", postHandler)

	server := httptest.NewServer(router)
	defer server.Close()

	t.Run("when GET method", func(t *testing.T) {
		response, err := http.Get(server.URL)
		if err != nil {
			t.Error("should not return error")
		}

		body, _ := ioutil.ReadAll(response.Body)
		response.Body.Close()

		if bodyString := string(body); bodyString != "handle GET method" {
			t.Errorf("should return the right response, %s", bodyString)
		}
	})

	t.Run("when POST method", func(t *testing.T) {
		response, err := http.Post(server.URL, "text/plain", nil)
		if err != nil {
			t.Error("should not return error")
		}

		body, _ := ioutil.ReadAll(response.Body)
		response.Body.Close()

		if bodyString := string(body); bodyString != "handle POST method" {
			t.Errorf("should return the right response, %s", bodyString)
		}
	})

	t.Run("when route does not exists", func(t *testing.T) {
		router := newRouter()
		server := httptest.NewServer(router)
		defer server.Close()

		response, err := http.Get(server.URL)
		if err != nil {
			t.Error("should not return error")
		}

		body, _ := ioutil.ReadAll(response.Body)
		response.Body.Close()

		if bodyString := string(body); bodyString != "No route matches" {
			t.Errorf("should return the right response, %s", bodyString)
		}
	})
}
