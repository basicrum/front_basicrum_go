package service

import (
	"encoding/json"
	"log"
	"time"

	"github.com/basicrum/front_basicrum_go/beacon"
	"github.com/basicrum/front_basicrum_go/dao"
	"github.com/basicrum/front_basicrum_go/geoip"
	"github.com/basicrum/front_basicrum_go/types"
	"github.com/ua-parser/uap-go/uaparser"
)

const hostUpdateDuration = time.Minute

// Service processes events and stores them in database access object
type Service struct {
	daoService      *dao.DAO
	userAgentParser *uaparser.Parser
	events          chan *types.Event
	geoIPService    geoip.Service
	hosts           map[string]string
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
		hosts:           map[string]string{},
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
	updateHostTicker := time.NewTicker(hostUpdateDuration)
	for {
		select {
		case event := <-s.events:
			s.processEvent(event)
		case <-updateHostTicker.C:
			s.processHosts()
		}
	}
}

func (s *Service) processEvent(event *types.Event) {
	if event == nil {
		return
	}
	beaconEvent := beacon.FromEvent(event)
	rumEvent := beacon.ConvertToRumEvent(beaconEvent, event, s.userAgentParser, s.geoIPService)
	jsonValue, err := json.Marshal(rumEvent)
	if err != nil {
		log.Printf("json parsing error: %+v", err)
		return
	}
	data := string(jsonValue)
	err = s.daoService.Save(data)
	if err != nil {
		log.Printf("failed to save data: %v err: %+v", data, err)
	}
	s.hosts[rumEvent.Hostname] = rumEvent.Created_At
}

func (s *Service) processHosts() {
	for hostname, createdAt := range s.hosts {
		s.saveHost(hostname, createdAt)
	}
	s.clearHosts()
}

func (s *Service) saveHost(hostname string, createdAt string) {
	event := beacon.NewHostnameEvent(hostname, createdAt)
	err := s.daoService.SaveHost(event)
	if err != nil {
		log.Printf("failed to save host: %+v err: %v", event, err)
	}
}

func (s *Service) clearHosts() {
	s.hosts = map[string]string{}
}
