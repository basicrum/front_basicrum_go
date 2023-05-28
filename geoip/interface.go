package geoip

// Service interface for geo ip
type Service interface {
	// CountryAndCity return country and city by http headers and remote ip address
	// nolint: revive
	CountryAndCity(ipString string) (string, string, error)
}
