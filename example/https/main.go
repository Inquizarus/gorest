package main

import (
	"net/http"
	"os"

	"github.com/inquizarus/gorest"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)
	log.SetOutput(os.Stdout)
	tlsConfig := gorest.ServeTLSConfig{
		Enabled:  true,
		CertPath: "restexample.crt",
		KeyPath:  "restexample.key",
	}
	config := gorest.ServeConfig{
		Port:   "8080",
		Logger: log,
		TLS:    tlsConfig,
		Handlers: []gorest.Handler{
			&gorest.BaseHandler{
				Path: "/",
				Get: func(w http.ResponseWriter, r *http.Request, params map[string]string) {
					defer r.Body.Close()
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("Hello, World!"))
				},
			},
		},
	}
	gorest.Serve(config)
}
