
package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_helloWorld(t *testing.T) {
	req, err := http.NewRequest("GET", "http://domain.com/", nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	helloWorld(res, req)

	exp := "Hello World"
	act := res.Body.String()
	if exp != act {
		t.Fatalf("Expected %s got %s", exp, act)
	}
}

