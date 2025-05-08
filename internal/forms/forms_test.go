package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	request := httptest.NewRequest("POST", "/test", nil)
	form := New(request.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("Invalid form: should have been valid")
	}
}

func TestForm_New(t *testing.T) {
	request := httptest.NewRequest("POST", "/test", nil)
	form := New(request.PostForm)

	if reflect.TypeOf(form) != reflect.TypeOf(&Form{}) {
		t.Errorf("Invalid type: expected %T, got %T", Form{}, form)
	}
}

func TestForm_Required(t *testing.T) {
	request := httptest.NewRequest("POST", "/test", nil)
	form := New(request.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("Invalid form: should be invalid - does not contain required fields")
	}

	postData := url.Values{}
	postData.Add("a", "a")
	postData.Add("b", "c")
	postData.Add("c", "a")

	request, _ = http.NewRequest("POST", "/test", nil)
	request.PostForm = postData
	form = New(request.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("Invalid form: should contain required fields")
	}
}

func TestForm_Has(t *testing.T) {
	request := httptest.NewRequest("POST", "/test", nil)
	form := New(request.PostForm)

	if form.Has("a") {
		t.Error("Invalid form: should be invalid - does not contain required field")
	}

	postData := url.Values{}
	postData.Add("a", "b")
	request, _ = http.NewRequest("POST", "/test", nil)
	request.PostForm = postData
	form = New(request.PostForm)
	if !form.Has("a") {
		t.Error("Invalid form: should contain required field")
	}
}

func TestForm_MinLength(t *testing.T) {
	request := httptest.NewRequest("POST", "/test", nil)
	form := New(request.PostForm)

	testString := "aaa"

	if !form.MinLength(testString, 0) {
		t.Error("Invalid form: length should be no bigger than 0")
	}

	postData := url.Values{}
	postData.Add("a", testString)
	request, _ = http.NewRequest("POST", "/test", nil)
	request.PostForm = postData
	form = New(request.PostForm)
	if !form.MinLength("a", len(testString)) {
		t.Errorf("Invalid form: attribute should be no bigger than %d", len(testString))
	}
}

func TestForm_IsEmail(t *testing.T) {
	postData := url.Values{}
	postData.Add("email", "razvan@gmailcom")
	postData.Add("email1", "razvangmail.com")
	postData.Add("email2", "razvan@gmail.com")
	request, _ := http.NewRequest("POST", "/test", nil)
	request.PostForm = postData
	form := New(request.PostForm)
	if form.IsEmail("email") {
		t.Error("Invalid form: attribute should not be valid email")
	}
	if form.IsEmail("email1") {
		t.Error("Invalid form: attribute should not be valid email")
	}
	if !form.IsEmail("email2") {
		t.Error("Invalid form: attribute should be valid email")
	}
}
