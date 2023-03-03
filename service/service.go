package service

import (
	"encoding/json"
	"log"

	"github.com/basicrum/front_basicrum_go/beacon"
	"github.com/basicrum/front_basicrum_go/dao"
	"github.com/basicrum/front_basicrum_go/types"
	"github.com/ua-parser/uap-go/uaparser"
)

// Service processes events and stores them in database access object
type Service struct {
	daoService      *dao.DAO
	userAgentParser *uaparser.Parser
	events          chan *types.Event
}

// New creates processing service
func New(
	daoService *dao.DAO,
	userAgentParser *uaparser.Parser,
) *Service {
	events := make(chan *types.Event)
	return &Service{
		daoService:      daoService,
		userAgentParser: userAgentParser,
		events:          events,
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
		rumEvent := beacon.ConvertToRumEvent(beaconEvent, s.userAgentParser)
		jsonValue, err := json.Marshal(rumEvent)
		if err != nil {
			log.Printf("json parsing error: %+v", err)
			continue
		}
		data := string(jsonValue)
		err = s.daoService.Save(data)
		if err != nil {
			log.Printf("failed to save data: %v err: %+v", data, err)
		}
	}
}
