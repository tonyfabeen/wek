package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	server := httptest.NewServer(handler{})
	defer server.Close()

	response, err := http.Get(server.URL)
	if err != nil {
		t.Error("should not return error")
	}

	body, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()

	if string(body) != "Hi from the outer space" {
		t.Error("should return the right response")
	}
}
