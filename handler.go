package gorest

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Handler interface for controllers
type Handler interface {
	GetPrefix() string
	GetPath() string
	Handle(http.ResponseWriter, *http.Request)
}

// VerbHandler is responseible for a specific type of request
type VerbHandler func(http.ResponseWriter, *http.Request, map[string]string)

// BaseHandler to use for extending
type BaseHandler struct {
	Path        string
	Prefix      string
	Name        string
	Middlewares []Middleware
	Get         VerbHandler
	Put         VerbHandler
	Post        VerbHandler
	Delete      VerbHandler
}

// GetPath for BaseHandler
func (h *BaseHandler) GetPath() string {
	return h.Path
}

// GetPrefix for BaseHandler
func (h *BaseHandler) GetPrefix() string {
	return h.Prefix
}

// Handle for BaseHandler
func (h *BaseHandler) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		executeHTTPFunc(h.Get, w, r, h.Middlewares)
		return
	case http.MethodPut:
		executeHTTPFunc(h.Put, w, r, h.Middlewares)
		return
	case http.MethodPost:
		executeHTTPFunc(h.Post, w, r, h.Middlewares)
		return
	case http.MethodDelete:
		executeHTTPFunc(h.Delete, w, r, h.Middlewares)
		return
	}
	executeHTTPFunc(nil, w, r, h.Middlewares)
}

func (h *BaseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Handle(w, r)
}

func executeHTTPFunc(f VerbHandler, w http.ResponseWriter, r *http.Request, mws []Middleware) {
	if nil == f {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	fn := func(w http.ResponseWriter, r *http.Request) {
		f(w, r, mux.Vars(r))
	}
	h := ChainMiddleware(http.HandlerFunc(fn), mws...)
	h.ServeHTTP(w, r)
}
