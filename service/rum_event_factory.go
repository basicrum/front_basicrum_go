package service

//go:generate mockgen -source=${GOFILE} -destination=mocks/${GOFILE} -package=servicemocks

import (
	"github.com/basicrum/front_basicrum_go/beacon"
	"github.com/basicrum/front_basicrum_go/geoip"
	"github.com/basicrum/front_basicrum_go/types"
	"github.com/ua-parser/uap-go/uaparser"
)

// IRumEventFactory rum event factory interface
type IRumEventFactory interface {
	// Create rum event from http captured event
	Create(event *types.Event) beacon.RumEvent
}

// RumEventFactory creates rum event
type RumEventFactory struct {
	userAgentParser *uaparser.Parser
	geoIPService    geoip.Service
}

// NewRumEventFactory creates rum event factory
func NewRumEventFactory(
	userAgentParser *uaparser.Parser,
	geoIPService geoip.Service,
) *RumEventFactory {
	return &RumEventFactory{
		userAgentParser: userAgentParser,
		geoIPService:    geoIPService,
	}
}

// Create rum event from http captured event
func (s *RumEventFactory) Create(event *types.Event) beacon.RumEvent {
	beaconEvent := beacon.FromEvent(event)
	return beacon.ConvertToRumEvent(beaconEvent, event, s.userAgentParser, s.geoIPService)
}
