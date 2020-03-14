package rest

import (
	"crypto/tls"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
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
