package gorest

import (
	log "github.com/sirupsen/logrus"
)

// ServeTLSConfig for serving HTTPS requests
type ServeTLSConfig struct {
	Enabled          bool
	Strict           bool
	StrictSubDomains bool
	Preload          bool
	CertPath         string
	KeyPath          string
}

// ServeConfig contains everything needed to serve the rest API
type ServeConfig struct {
	Handlers    []Handler
	Middlewares []Middleware
	Port        string
	TLS         ServeTLSConfig
	Logger      log.StdLogger
}

// GetPort returns port from config if set else check
// if TLS is enabled and returns 443 in that case and lastly fallback
// to default HTTP port 80
func (sc *ServeConfig) GetPort() string {
	if "" != sc.Port {
		return sc.Port
	}
	if true == sc.TLS.Enabled {
		return "443"
	}
	return "80"
}
