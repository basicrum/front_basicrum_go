package service

import (
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
		Create    bool
		Save      bool
		SaveError error
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
			name: "should save the event into database",
			expects: expects{
				Create: true,
				Save:   true,
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
			backupService := backupmocks.NewMockIBackup(ctrl)

			s := New(
				rumEventFactory,
				daoService,
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

			// expects
			rumEvent := beacon.RumEvent{
				Hostname: hostname,
			}
			if tt.expects.Create {
				rumEventFactory.EXPECT().Create(testEvent).Return(rumEvent)
			}
			if tt.expects.Save {
				daoService.EXPECT().Save(rumEvent).Return(tt.expects.SaveError)
			}

			// when
			s.processEvent(inputEvent)
		})
	}
}
