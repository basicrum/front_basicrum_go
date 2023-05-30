package cloudflare

import (
	"net/http"
	"strings"
)

// Service implement cloudflare geoip service
type Service struct {
}

// New creates a new service
func New() *Service {
	return &Service{}
}

// CountryAndCity return country and city by http headers and remote ip address
// nolint: revive
func (s *Service) CountryAndCity(header http.Header, _ string) (string, string, error) {
	country := cleanupHeaderValue(header.Get("CF-IPCountry"))
	city := cleanupHeaderValue(header.Get("CF-IPCity"))
	return country, city, nil
}

func cleanupHeaderValue(hVal string) string {
	hVal = strings.TrimSpace(hVal)
	hVal = strings.TrimPrefix(hVal, "\"")
	hVal = strings.TrimSuffix(hVal, "\"")
	hVal = strings.TrimSpace(hVal)
	return hVal
}
