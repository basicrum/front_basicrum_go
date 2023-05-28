package maxmind

import (
	"fmt"
	"net"

	_ "embed"

	"github.com/oschwald/geoip2-golang"
)

//go:embed GeoLite2-City.mmdb
var geoLite2City []byte

// Service implement maxmind geoip service
type Service struct {
}

// New creates a new service
func New() *Service {
	return &Service{}
}

// CountryAndCity return country and city by http headers and remote ip address
// nolint: revive
func (s *Service) CountryAndCity(ipString string) (string, string, error) {
	db, err := geoip2.FromBytes(geoLite2City)
	if err != nil {
		return "", "", err
	}
	defer db.Close()

	ip := net.ParseIP(ipString)
	if ip == nil {
		return "", "", fmt.Errorf("cannot parse ip[%v]", ipString)
	}

	record, err := db.City(ip)
	if err != nil {
		return "", "", err
	}

	return record.Country.IsoCode, record.City.Names["en"], nil
}
