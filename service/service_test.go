package service

import (
	"testing"

	"github.com/basicrum/front_basicrum_go/backup"
	"github.com/basicrum/front_basicrum_go/dao"
	"github.com/basicrum/front_basicrum_go/geoip"
	"github.com/basicrum/front_basicrum_go/types"
	"github.com/ua-parser/uap-go/uaparser"
)

func TestService_processEvent(t *testing.T) {
	type fields struct {
		daoService          *dao.DAO
		userAgentParser     *uaparser.Parser
		events              chan *types.Event
		geoIPService        geoip.Service
		hosts               map[string]string
		subscriptionService ISubscriptionService
		backupService       backup.IBackup
	}
	type args struct {
		event *types.Event
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				daoService:          tt.fields.daoService,
				userAgentParser:     tt.fields.userAgentParser,
				events:              tt.fields.events,
				geoIPService:        tt.fields.geoIPService,
				hosts:               tt.fields.hosts,
				subscriptionService: tt.fields.subscriptionService,
				backupService:       tt.fields.backupService,
			}
			s.processEvent(tt.args.event)
		})
	}
}
