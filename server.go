package rest

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Serve the rest service
func Serve(config ServeConfig) {
	r := mux.NewRouter()
	registerHandlers(r, config)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.GetPort()),
		Handler: r,
	}

	if true == config.TLS.Enabled {
		server.TLSConfig = generateTLSConfig()
		server.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0)
		config.Logger.Printf("server started with TLS enabled on port %s", config.GetPort())
		config.Logger.Println(server.ListenAndServeTLS(
			config.TLS.CertPath,
			config.TLS.KeyPath,
		))
	}
	if false == config.TLS.Enabled {
		config.Logger.Printf("server started on port %s", config.GetPort())
		config.Logger.Println(server.ListenAndServe())
	}
}

func generateTLSConfig() *tls.Config {
	return &tls.Config{
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
}

func registerHandlers(r *mux.Router, config ServeConfig) {
	for _, h := range config.Handlers {
		r.HandleFunc(h.GetPath(), ChainMiddleware(h.Handle, config.Middlewares...)).Name(h.GetName())
	}
}
