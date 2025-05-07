package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var testHandler myHTTPHandler
	h := NoSurf(&testHandler)

	// check what the return type of NoSurf is
	switch v := h.(type) {
	case http.Handler:
		// nothing, pass
	default:
		// return type should be http.Handler
		t.Error(fmt.Sprintf("NoSurf() return type %T is not http.Handler", v))
	}
}

func TestSessionLoad(t *testing.T) {
	var testHandler myHTTPHandler
	h := SessionLoad(&testHandler)

	// check what the return type of NoSurf is
	switch v := h.(type) {
	case http.Handler:
		// nothing, pass
	default:
		// return type should be http.Handler
		t.Error(fmt.Sprintf("SessionLoad() return type %T is not http.Handler", v))
	}
}
