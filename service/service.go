package service

//go:generate mockgen -source=${GOFILE} -destination=mocks/${GOFILE} -package=servicemocks

import (
	"encoding/json"
	"log"
	"time"

	"github.com/basicrum/front_basicrum_go/backup"
	"github.com/basicrum/front_basicrum_go/beacon"
	"github.com/basicrum/front_basicrum_go/dao"
	"github.com/basicrum/front_basicrum_go/geoip"
	"github.com/basicrum/front_basicrum_go/types"
	"github.com/ua-parser/uap-go/uaparser"
)

const hostUpdateDuration = time.Minute

// IService service interface
type IService interface {
	// Run runs the service
	Run()
	// SaveAsync saves an event asynchronously
	SaveAsync(event *types.Event)
	// RegisterHostname generates new subscription
	RegisterHostname(hostname, username string) error
	// DeleteHostname deletes the hostname
	DeleteHostname(hostname, username string) error
}

// Service processes events and stores them in database access object
type Service struct {
	daoService          dao.IDAO
	userAgentParser     *uaparser.Parser
	events              chan *types.Event
	geoIPService        geoip.Service
	hosts               map[string]string
	subscriptionService ISubscriptionService
	backupService       backup.IBackup
}

// New creates processing service
// nolint: revive
func New(
	daoService dao.IDAO,
	userAgentParser *uaparser.Parser,
	geoIPService geoip.Service,
	subscriptionService ISubscriptionService,
	backupService backup.IBackup,
) *Service {
	events := make(chan *types.Event)
	return &Service{
		daoService:          daoService,
		userAgentParser:     userAgentParser,
		events:              events,
		geoIPService:        geoIPService,
		hosts:               map[string]string{},
		subscriptionService: subscriptionService,
		backupService:       backupService,
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

// RegisterHostname generates new subscription
func (s *Service) RegisterHostname(hostname, username string) error {
	subscription := types.NewSubscription(time.Now())
	ownerHostname := types.NewOwnerHostname(username, hostname, subscription)
	return s.daoService.InsertOwnerHostname(ownerHostname)
}

// DeleteHostname deletes the hostname
func (s *Service) DeleteHostname(hostname, username string) error {
	return s.daoService.DeleteOwnerHostname(hostname, username)
}

func (s *Service) processEvent(event *types.Event) {
	if event == nil {
		return
	}
	beaconEvent := beacon.FromEvent(event)
	rumEvent := beacon.ConvertToRumEvent(beaconEvent, event, s.userAgentParser, s.geoIPService)
	lookup, err := s.subscriptionService.GetSubscription(rumEvent.SubscriptionID, rumEvent.Hostname)
	if err != nil {
		log.Printf("get subscription error: %+v", err)
		return
	}

	switch lookup {
	case types.NewFoundLookup().Value:
		s.processRumEvent(rumEvent)
	case types.NewExpiredLookup().Value:
		s.backupService.SaveExpired(event)
	case types.NewNotFoundLookup().Value:
		s.backupService.SaveUnknown(event)
	default:
		log.Printf("unsupported lookup result: %v", lookup)
		return
	}
}

func (s *Service) processRumEvent(rumEvent beacon.RumEvent) {
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
