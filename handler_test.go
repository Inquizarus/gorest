package gorest_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	rest "github.com/inquizarus/gorest"
)

func TestThatBaseHandlerGetPathReturnsCorrectValue(t *testing.T) {
	expected := "/"
	handler := rest.BaseHandler{
		Path: expected,
	}
	actual := handler.GetPath()
	if !assertEquals(actual, expected) {
		t.Errorf("handler GetPath return wrong value, expected %s but got %s", expected, actual)
	}
}

func TestThatHandleDirectsGETRequestsCorrectly(t *testing.T) {
	cases := []string{"GET", "PUT", "POST", "DELETE"}
	f := func(w http.ResponseWriter, r *http.Request, vars map[string]string) {
		w.Write([]byte("called"))
	}
	h := rest.BaseHandler{
		Get:    f,
		Put:    f,
		Post:   f,
		Delete: f,
	}

	for _, c := range cases {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(c, "/", strings.NewReader(""))
		h.Handle(w, r)
		bc, _ := ioutil.ReadAll(w.Body)
		if !assertEquals(string(bc), "called") {
			t.Errorf("%s func where never called when using Handle function in handler for %s request\n", c, c)
		}
		h.ServeHTTP(w, r)
		bc, _ = ioutil.ReadAll(w.Body)
		if !assertEquals(string(bc), "called") {
			t.Errorf("%s func where never called when using ServeHTTP function in handler for %s request\n", c, c)
		}
	}
}

func TestThatNotFoundIsTriggeredForUnSupportedRequestMethods(t *testing.T) {
	h := rest.BaseHandler{}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("UNSUPPORTED", "/", strings.NewReader(""))
	h.Handle(w, r)
	if w.Code != http.StatusNotFound {
		fmt.Println("response headers did not have \"not found\" status for unsupported request method")
	}
}

func TestThatBaseHandlerGetPrefixReturnsCorrectValue(t *testing.T) {
	expected := "/test/"
	handler := rest.BaseHandler{
		Prefix: expected,
	}
	actual := handler.GetPrefix()
	if !assertEquals(actual, expected) {
		t.Errorf("handler GetPrefix return wrong value, expected %s but got %s", expected, actual)
	}
}

func assertEquals(x interface{}, y interface{}) bool {
	return x == y
}
