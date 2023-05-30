package geoip

import (
	"net/http"
)

// Composite chain two geoip services
type Composite struct {
	primary Service
	next    Service
}

// NewComposite creates a composite geoip service
func NewComposite(primary Service, next Service) *Composite {
	return &Composite{
		primary: primary,
		next:    next,
	}
}

// CountryAndCity return country and city by http headers and remote ip address
// nolint: revive
func (s *Composite) CountryAndCity(header http.Header, ipString string) (string, string, error) {
	country, city, err := s.primary.CountryAndCity(header, ipString)
	if (err != nil) || (country == "" && city == "") {
		return s.next.CountryAndCity(header, ipString)
	}
	return country, city, err
}
