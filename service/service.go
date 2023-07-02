package service

import (
	"log"

	"github.com/basicrum/front_basicrum_go/beacon"
	"github.com/basicrum/front_basicrum_go/dao"
	"github.com/basicrum/front_basicrum_go/geoip"
	"github.com/basicrum/front_basicrum_go/types"
	"github.com/ua-parser/uap-go/uaparser"
)

// Service processes events and stores them in database access object
type Service struct {
	daoService      *dao.DAO
	userAgentParser *uaparser.Parser
	events          chan *types.Event
	geoIPService    geoip.Service
}

// New creates processing service
func New(
	daoService *dao.DAO,
	userAgentParser *uaparser.Parser,
	geoIPService geoip.Service,
) *Service {
	events := make(chan *types.Event)
	return &Service{
		daoService:      daoService,
		userAgentParser: userAgentParser,
		events:          events,
		geoIPService:    geoIPService,
	}
}

// SaveAsync saves an event asynchronously
func (s *Service) SaveAsync(event *types.Event) {
	go func() {
		s.events <- event
	}()
}

// Run process the events from the channel and save them in datastore (click house)
func (s *Service) Run() {
	for {
		event := <-s.events
		if event == nil {
			continue
		}
		beaconEvent := beacon.FromEvent(event)
		rumEvent := beacon.ConvertToRumEvent(beaconEvent, event, s.userAgentParser, s.geoIPService)
		err := s.daoService.Save(rumEvent)
		if err != nil {
			log.Printf("failed to save data: %+v err: %+v", rumEvent, err)
		}
	}
}
