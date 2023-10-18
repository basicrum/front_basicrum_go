package service

import (
	"bufio"
	"fmt"
	"io"
	"net/url"
	"os"
	"testing"

	backupmocks "github.com/basicrum/front_basicrum_go/backup/mocks"
	"github.com/basicrum/front_basicrum_go/beacon"
	daomocks "github.com/basicrum/front_basicrum_go/dao/mocks"
	"github.com/basicrum/front_basicrum_go/geoip"
	"github.com/basicrum/front_basicrum_go/geoip/cloudflare"
	"github.com/basicrum/front_basicrum_go/geoip/maxmind"
	servicemocks "github.com/basicrum/front_basicrum_go/service/mocks"
	"github.com/basicrum/front_basicrum_go/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/ua-parser/uap-go/uaparser"
)

func TestService_processEvent(t *testing.T) {
	type args struct {
		event          *types.Event
		subscriptionID string
		hostname       string
	}
	type expects struct {
		GetSubscription      bool
		lookup               *types.Lookup
		GetSubscriptionError error
		processRumEvent      bool
		rumEvent             beacon.RumEvent
		SaveExpired          bool
		SaveExpiredEvent     *types.Event
		SaveUnknown          bool
		SaveUnknownEvent     *types.Event
	}
	tests := []struct {
		name    string
		args    args
		expects expects
		want    string
		wantErr bool
	}{
		{
			name: "Get subscription lookup",
			args: args{
				subscriptionID: "subscription_id1",
				hostname:       "host1",
			},
			expects: expects{
				GetSubscription: true,
				lookup:          types.NewFoundLookup().Value,
			},
			want:    string(types.FoundLookup),
			wantErr: false,
		},
		{
			name: "Process found event",
			args: args{
				event: &types.Event{
					RequestParameters: url.Values{
						"hostname":        []string{"hostname1"},
						"subscription_id": []string{"subscription_id1"},
						"created_at":      []string{"created_at1"},
						"user_agent":      []string{"Chrome/104.0.5112.102"},
					},
				},
			},
			expects: expects{
				processRumEvent: true,
				rumEvent:        beacon.RumEvent{},
			},
			want: string(*types.NewFoundLookup().Value),
		},
		{
			name: "Lookup subscription expired",
			args: args{
				event: &types.Event{
					RequestParameters: url.Values{
						"hostname":        []string{"hostname1"},
						"subscription_id": []string{"subscription_id1"},
						"created_at":      []string{"created_at1"},
						"user_agent":      []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.5112.102 Safari/537.36"},
					},
				},
			},
			expects: expects{
				SaveExpired: true,
				SaveExpiredEvent: &types.Event{
					RequestParameters: url.Values{
						"hostname":        []string{"hostname1"},
						"subscription_id": []string{"subscription_id1"},
						"created_at":      []string{"created_at1"},
						"user_agent":      []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.5112.102 Safari/537.36"},
					},
				},
			},
			want: string(*types.NewExpiredLookup().Value),
		},
		{
			name: "Lookup subscription not found",
			args: args{
				event: &types.Event{},
			},
			expects: expects{
				SaveUnknown: true,
				SaveUnknownEvent: &types.Event{
					RequestParameters: url.Values{
						"hostname":        []string{"hostname1"},
						"subscription_id": []string{"subscription_id1"},
						"created_at":      []string{"created_at1"},
						"user_agent":      []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.5112.102 Safari/537.36"},
					},
				},
			},
			want: string(*types.NewNotFoundLookup().Value),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			daoService := daomocks.NewMockIDAO(ctrl)
			file, err := os.Open("../assets/uaparser_regexes.yaml")
			require.NoError(t, err)
			defer file.Close()
			// Get the file size
			stat, err := file.Stat()
			if err != nil {
				fmt.Println(err)
				return
			}
			// Read the file into a byte slice
			userAgentRegularExpressions := make([]byte, stat.Size())
			_, err = bufio.NewReader(file).Read(userAgentRegularExpressions)
			if err != nil && err != io.EOF {
				fmt.Println(err)
				return
			}
			userAgentParser, err := uaparser.NewFromBytes(userAgentRegularExpressions)
			require.NoError(t, err)
			geopIPService := geoip.NewComposite(
				cloudflare.New(),
				maxmind.New(),
			)
			backupService := backupmocks.NewMockIBackup(ctrl)
			subscriptionService := servicemocks.NewMockISubscriptionService(ctrl)
			s := New(daoService, userAgentParser, geopIPService, subscriptionService, backupService)
			if tt.expects.GetSubscription {
				subscriptionService.EXPECT().
					GetSubscription(
						tt.expects.rumEvent.SubscriptionID,
						tt.expects.rumEvent.Hostname,
					).Return(tt.expects.lookup, tt.expects.GetSubscriptionError)
			}
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
			var got *types.Lookup
			got, err = s.subscriptionService.GetSubscription(tt.args.subscriptionID, tt.args.hostname)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.want, got)
			s.processEvent(tt.args.event)
		})
	}
}
