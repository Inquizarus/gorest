package gorest

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestThatHandlersGetRegistered(t *testing.T) {

	r := mux.NewRouter()

	config := ServeConfig{
		Handlers: []Handler{
			&BaseHandler{
				Name: "test-route",
				Path: "/",
			},
			&BaseHandler{
				Name: "test-route-sub",
				Path: "/sub",
			},
		},
	}

	registerHandlers(r, config)

	registered := []string{}

	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		registered = append(registered, route.GetName())
		return nil
	})

	for _, h := range config.Handlers {
		found := false
		for _, n := range registered {
			if n == h.GetName() {
				found = true
			}
		}
		if false == found {
			t.Errorf("route %s not registered", h.GetPath())
		}
	}

}

func TestThatGenerateTLSConfigWorks(t *testing.T) {
	actual := generateTLSConfig()
	expected := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}
	if true != reflect.DeepEqual(actual, expected) {
		t.Error("generated TLS configuration does not correspond to expected one, check this test for correct configuration")
	}
}

func TestThatTheCorrectHandlerIsServed(t *testing.T) {

	r := mux.NewRouter()

	config := ServeConfig{
		Handlers: []Handler{
			&BaseHandler{
				Name: "test-route",
				Path: "/",
				Get: func(w http.ResponseWriter, r *http.Request, p map[string]string) {
					defer r.Body.Close()
					w.Write([]byte("root path"))
				},
				Post: func(w http.ResponseWriter, r *http.Request, p map[string]string) {
					defer r.Body.Close()
					w.Write([]byte("post root path"))
				},
			},
			&BaseHandler{
				Name: "test-route-sub",
				Path: "/sub",
				Get: func(w http.ResponseWriter, r *http.Request, p map[string]string) {
					defer r.Body.Close()
					w.Write([]byte("sub path"))
				},
			},
			&BaseHandler{
				Name: "test-route-sub-parameterized",
				Path: "/sub/{name}",
				Get: func(w http.ResponseWriter, r *http.Request, p map[string]string) {
					defer r.Body.Close()
					name := p["name"]
					w.Write([]byte(fmt.Sprintf("%s sub path", name)))
				},
			},
		},
	}

	registerHandlers(r, config)

	cases := []struct {
		Path     string
		Method   string
		Expected string
	}{
		{Path: "/", Method: "GET", Expected: "root path"},
		{Path: "/", Method: "POST", Expected: "post root path"},
		{Path: "/sub", Method: "GET", Expected: "sub path"},
		{Path: "/sub/parameterized", Method: "GET", Expected: "parameterized sub path"},
	}

	for _, c := range cases {
		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, httptest.NewRequest(c.Method, c.Path, nil))
		body := recorder.Body.String()
		assert.Equal(t, c.Expected, body, fmt.Sprintf("expected %s but got %s", c.Expected, body))
	}
}
