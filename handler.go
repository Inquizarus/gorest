package rest

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Handler interface for controllers
type Handler interface {
	GetPath() string
	GetName() string
	Handle(http.ResponseWriter, *http.Request)
}

// BaseHandler to use for extending
type BaseHandler struct {
	Path        string
	Name        string
	Middlewares []Middleware
	Get         func(http.ResponseWriter, *http.Request, map[string]string)
	Put         func(http.ResponseWriter, *http.Request, map[string]string)
	Post        func(http.ResponseWriter, *http.Request, map[string]string)
	Delete      func(http.ResponseWriter, *http.Request, map[string]string)
}

// GetPath for BaseHandler
func (h *BaseHandler) GetPath() string {
	return h.Path
}

// GetName for BaseHandler
func (h *BaseHandler) GetName() string {
	return h.Name
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

func executeHTTPFunc(f func(http.ResponseWriter, *http.Request, map[string]string), w http.ResponseWriter, r *http.Request, mws []Middleware) {
	if nil == f {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	ChainMiddleware(func(w http.ResponseWriter, r *http.Request) {
		f(w, r, mux.Vars(r))
	}, mws...)(w, r)
}
