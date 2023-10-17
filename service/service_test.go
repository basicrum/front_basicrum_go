package service

import (
	"testing"

	backupmocks "github.com/basicrum/front_basicrum_go/backup/mocks"
	"github.com/basicrum/front_basicrum_go/beacon"
	"github.com/basicrum/front_basicrum_go/dao"
	"github.com/basicrum/front_basicrum_go/geoip"
	"github.com/basicrum/front_basicrum_go/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ua-parser/uap-go/uaparser"
)

func TestService_processEvent(t *testing.T) {
	type fields struct {
		daoService          *dao.DAO
		userAgentParser     *uaparser.Parser
		events              chan *types.Event
		geoIPService        geoip.Service
		subscriptionService ISubscriptionService
	}
	type args struct {
		rumEvent beacon.RumEvent
		lookup   Lookup
	}
	type expects struct {
		GetSubscription  bool
		processRumEvent  bool
		rumEvent         beacon.RumEvent
		SaveExpired      bool
		SaveExpiredEvent *types.Event
		SaveUnknown      bool
		SaveUnknownEvent *types.Event
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		expects expects
		want    string
	}{
		{
			name: "Lookup subscription",
			args: args{
				rumEvent: beacon.RumEvent{},
			},
			expects: expects{
				GetSubscription: true,
				rumEvent:        beacon.RumEvent{},
			},
			want: string(FoundLookup),
		},
		{
			name: "Lookup subscription found",
			args: args{
				lookup: FoundLookup,
			},
			expects: expects{
				processRumEvent: true,
				rumEvent:        beacon.RumEvent{},
			},
			want: string(ExpiredLookup),
		},
		{
			name: "Lookup subscription expired",
			args: args{
				lookup: FoundLookup,
			},
			expects: expects{
				SaveExpired:      true,
				SaveExpiredEvent: &types.Event{},
			},
			want: string(ExpiredLookup),
		},
		{
			name: "Lookup subscription not found",
			args: args{
				lookup: ExpiredLookup,
			},
			expects: expects{
				SaveUnknown:      true,
				SaveUnknownEvent: &types.Event{},
			},
			want: string(NotFoundLookup),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// s := &Service{
			// 	daoService:          tt.fields.daoService,
			// 	userAgentParser:     tt.fields.userAgentParser,
			// 	events:              tt.fields.events,
			// 	geoIPService:        tt.fields.geoIPService,
			// 	hosts:               tt.fields.hosts,
			// 	subscriptionService: tt.fields.subscriptionService,
			// 	backupService:       tt.fields.backupService,
			// }
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			backupService := backupmocks.NewMockIBackup(ctrl)
			s := New(tt.fields.daoService, tt.fields.userAgentParser, tt.fields.geoIPService, tt.fields.subscriptionService, backupService)
			for e := range tt.fields.events {
				// if tt.expects.processRumEvent {
				// 	s.processRumEvent.EXPECT().
				// 		processRumEvent(
				// 			tt.expects.rumEvent,
				// 		)
				// }
				if tt.expects.SaveExpired {
					backupService.EXPECT().
						SaveExpired(
							tt.expects.SaveExpiredEvent,
						)
				}
				if tt.expects.SaveUnknown {
					backupService.EXPECT().
						SaveUnknown(
							tt.expects.SaveUnknownEvent,
						)
				}
				b := beacon.FromEvent(e)
				re := beacon.ConvertToRumEvent(b, e, tt.fields.userAgentParser, tt.fields.geoIPService)
				lk, err := s.subscriptionService.GetSubscription(re.SubscriptionID, re.Hostname)
				require.NoError(t, err)
				assert.Equal(t, tt.want, lk)
			}
		})
	}
}
