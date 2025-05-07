package main

import (
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	os.Exit(m.Run())
}

type myHTTPHandler struct {
}

func (mh *myHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
