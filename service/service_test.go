package service

import (
	"errors"
	"net/url"
	"testing"

	backupmocks "github.com/basicrum/front_basicrum_go/backup/mocks"
	"github.com/basicrum/front_basicrum_go/beacon"
	daomocks "github.com/basicrum/front_basicrum_go/dao/mocks"
	servicemocks "github.com/basicrum/front_basicrum_go/service/mocks"
	"github.com/basicrum/front_basicrum_go/types"
	"github.com/golang/mock/gomock"
)

func TestService_processEvent(t *testing.T) {
	type expects struct {
		Create                bool
		GetSubscription       bool
		GetSubscriptionReturn Lookup
		GetSubscriptionError  error
		Save                  bool
		SaveError             error
		SaveExpired           bool
		SaveUnknown           bool
	}
	type args struct {
		nilEvent bool
	}
	tests := []struct {
		name    string
		expects expects
		args    args
	}{
		{
			name: "when event is nil expect do nothing",
			args: args{
				nilEvent: true,
			},
		},
		{
			name: "when GetSubscription return FoundLookup should save the event into database",
			expects: expects{
				Create:                true,
				GetSubscription:       true,
				GetSubscriptionReturn: FoundLookup,
				Save:                  true,
			},
		},
		{
			name: "when GetSubscription return FoundLookup and save database return error should ignore the error",
			expects: expects{
				Create:                true,
				GetSubscription:       true,
				GetSubscriptionReturn: FoundLookup,
				Save:                  true,
				SaveError:             errors.New("error1"),
			},
		},
		{
			name: "when GetSubscription return ExpiredLookup should call backup SaveExpired",
			expects: expects{
				Create:                true,
				GetSubscription:       true,
				GetSubscriptionReturn: ExpiredLookup,
				SaveExpired:           true,
			},
		},
		{
			name: "when GetSubscription return NotFoundLookup should call backup SaveUnknown",
			expects: expects{
				Create:                true,
				GetSubscription:       true,
				GetSubscriptionReturn: NotFoundLookup,
				SaveUnknown:           true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// init
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			daoService := daomocks.NewMockIDAO(ctrl)
			rumEventFactory := servicemocks.NewMockIRumEventFactory(ctrl)
			subscriptionService := NewMockISubscriptionService(ctrl)
			backupService := backupmocks.NewMockIBackup(ctrl)

			s := New(
				rumEventFactory,
				daoService,
				subscriptionService,
				backupService,
			)

			// given
			testEvent := &types.Event{
				RequestParameters: url.Values{
					"key1": []string{"value1", "value2"},
				},
			}
			var inputEvent *types.Event
			if !tt.args.nilEvent {
				inputEvent = testEvent
			}
			hostname := "hostname1"
			subscriptionID := "subscriptionID1"

			// expects
			rumEvent := beacon.RumEvent{
				Hostname:       hostname,
				SubscriptionID: subscriptionID,
			}
			if tt.expects.Create {
				rumEventFactory.EXPECT().Create(testEvent).Return(rumEvent)
			}
			if tt.expects.GetSubscription {
				subscriptionService.EXPECT().GetSubscription(subscriptionID, hostname).Return(tt.expects.GetSubscriptionReturn, tt.expects.GetSubscriptionError)
			}
			if tt.expects.Save {
				daoService.EXPECT().Save(rumEvent).Return(tt.expects.SaveError)
			}
			if tt.expects.SaveExpired {
				backupService.EXPECT().SaveExpired(testEvent)
			}
			if tt.expects.SaveUnknown {
				backupService.EXPECT().SaveUnknown(testEvent)
			}

			// when
			s.processEvent(inputEvent)
		})
	}
}
