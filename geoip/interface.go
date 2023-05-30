package geoip

import "net/http"

// Service interface for geo ip
type Service interface {
	// CountryAndCity return country and city by http headers and remote ip address
	// nolint: revive
	CountryAndCity(header http.Header, ipString string) (string, string, error)
}
