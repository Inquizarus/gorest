package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestThatWithJSONContentSetTheCorrectContentType(t *testing.T) {
	mw := WithJSONContent()
	s := httptest.NewServer(mw(func(w http.ResponseWriter, r *http.Request) {}))
	defer s.Close()
	res, err := http.Get(s.URL)
	if err != nil {
		t.Fatal(err)
	}
	ct := res.Header.Get("content-type")
	res.Body.Close()
	if "application/json" != ct {
		t.Error("content-type where not the expected applicaton/json after middlewere ran")
	}
}

func TestThatWithStrictTransportSecuritySetCorrectHeader(t *testing.T) {
	cases := []struct {
		Expected string
		TLS      ServeTLSConfig
	}{
		{Expected: "", TLS: ServeTLSConfig{Strict: false}},
		{Expected: "max-age=63072000", TLS: ServeTLSConfig{Strict: true}},
		{Expected: "max-age=63072000; includeSubDomains", TLS: ServeTLSConfig{Strict: true, StrictSubDomains: true}},
		{Expected: "max-age=63072000; includeSubDomains; preload", TLS: ServeTLSConfig{Strict: true, StrictSubDomains: true, Preload: true}},
		{Expected: "max-age=63072000; preload", TLS: ServeTLSConfig{Strict: true, StrictSubDomains: false, Preload: true}},
	}
	for _, c := range cases {
		mw := WithStrictTransportSecurity(c.TLS)
		s := httptest.NewServer(mw(func(w http.ResponseWriter, r *http.Request) {}))
		defer s.Close()
		res, err := http.Get(s.URL)
		if err != nil {
			t.Fatal(err)
		}
		sts := res.Header.Get("Strict-Transport-Security")
		res.Body.Close()
		if c.Expected != sts {
			t.Errorf("Strict-Transport-Security header did not have the expected value, got %s but wanted %s", sts, c.Expected)
		}
	}
}

func TestThatChainMiddlewareRunsAllPassedMiddlewares(t *testing.T) {
	ran := false
	mw := func(f http.HandlerFunc) http.HandlerFunc {
		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {
			ran = true
			f(w, r)
		}
	}
	s := httptest.NewServer(ChainMiddleware(func(w http.ResponseWriter, r *http.Request) {}, mw))
	defer s.Close()
	_, err := http.Get(s.URL)
	if err != nil {
		t.Fatal(err)
	}
	if false == ran {
		t.Error("middlewares where not run propperly when chaining them")
	}
}
