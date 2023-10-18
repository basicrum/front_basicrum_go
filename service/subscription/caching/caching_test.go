package caching

import (
	"errors"
	"testing"
	"time"

	daomocks "github.com/basicrum/front_basicrum_go/dao/mocks"
	"github.com/basicrum/front_basicrum_go/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCacheSubscriptionService_GetSubscription(t *testing.T) {
	validTime := time.Now().Add(time.Hour)
	expiredTime := time.Now().Add(-time.Hour)
	type expects struct {
		GetSubscriptionsReturn map[string]*types.SubscriptionWithHostname
		GetSubscriptionsError  error
		GetSubscription        bool
		GetSubscriptionTimes   int
		GetSubscriptionRequest string
		GetSubscriptionReturn  *types.SubscriptionWithHostname
		GetSubscriptionError   error
	}
	type args struct {
		subscriptionID string
		hostname       string
	}
	tests := []struct {
		name                    string
		args                    args
		expects                 expects
		want                    *types.Lookup
		wantGetSubscriptionsErr bool
		wantErr                 bool
	}{
		{
			name: "when subscription id is found in the cache then expected Found",
			args: args{
				subscriptionID: "subscriptionID1",
				hostname:       "hostname1",
			},
			expects: expects{
				GetSubscriptionsReturn: map[string]*types.SubscriptionWithHostname{
					"subscriptionID1": {
						Subscription: types.Subscription{
							ID:        "subscriptionID1",
							ExpiresAt: validTime,
						},
						Hostname: "hostname1",
					},
				},
			},
			want: types.NewFoundLookup().Value,
		},
		{
			name: "when subscription id is found in the cache with different hostname then expected NotFound",
			args: args{
				subscriptionID: "subscriptionID1",
				hostname:       "hostname1",
			},
			expects: expects{
				GetSubscriptionsReturn: map[string]*types.SubscriptionWithHostname{
					"subscriptionID1": {
						Subscription: types.Subscription{
							ID:        "subscriptionID1",
							ExpiresAt: validTime,
						},
						Hostname: "otherHostname",
					},
				},
			},
			want: types.NewNotFoundLookup().Value,
		},
		{
			name: "when subscription id is not found in the cache then load from dao expected Found",
			args: args{
				subscriptionID: "subscriptionID1",
				hostname:       "hostname1",
			},
			expects: expects{
				GetSubscriptionsReturn: map[string]*types.SubscriptionWithHostname{},
				GetSubscription:        true,
				GetSubscriptionRequest: "subscriptionID1",
				GetSubscriptionReturn: &types.SubscriptionWithHostname{
					Subscription: types.Subscription{
						ID:        "subscriptionID1",
						ExpiresAt: validTime,
					},
					Hostname: "hostname1",
				},
			},
			want: types.NewFoundLookup().Value,
		},
		{
			name: "when subscription id is found expired in the cache then load from dao expected Expired",
			args: args{
				subscriptionID: "subscriptionID1",
				hostname:       "hostname1",
			},
			expects: expects{
				GetSubscriptionsReturn: map[string]*types.SubscriptionWithHostname{},
				GetSubscription:        true,
				GetSubscriptionRequest: "subscriptionID1",
				GetSubscriptionReturn: &types.SubscriptionWithHostname{
					Subscription: types.Subscription{
						ID:        "subscriptionID1",
						ExpiresAt: expiredTime,
					},
					Hostname: "hostname1",
				},
			},
			want: types.NewExpiredLookup().Value,
		},
		{
			name: "when subscription id is not found in the cache then load from dao with different hostname expected NotFound",
			args: args{
				subscriptionID: "subscriptionID1",
				hostname:       "hostname1",
			},
			expects: expects{
				GetSubscriptionsReturn: map[string]*types.SubscriptionWithHostname{},
				GetSubscription:        true,
				GetSubscriptionRequest: "subscriptionID1",
				GetSubscriptionReturn: &types.SubscriptionWithHostname{
					Subscription: types.Subscription{
						ID:        "subscriptionID1",
						ExpiresAt: validTime,
					},
					Hostname: "otherHostname",
				},
			},
			want: types.NewNotFoundLookup().Value,
		},
		{
			name: "when subscription id is not found in the cache then load from dao expired expected Expired",
			args: args{
				subscriptionID: "subscriptionID1",
				hostname:       "hostname1",
			},
			expects: expects{
				GetSubscriptionsReturn: map[string]*types.SubscriptionWithHostname{},
				GetSubscription:        true,
				GetSubscriptionRequest: "subscriptionID1",
				GetSubscriptionReturn: &types.SubscriptionWithHostname{
					Subscription: types.Subscription{
						ID:        "subscriptionID1",
						ExpiresAt: expiredTime,
					},
					Hostname: "hostname1",
				},
			},
			want: types.NewExpiredLookup().Value,
		},
		{
			name: "when subscription id is not found in the cache then not found in dao expected NotFound",
			args: args{
				subscriptionID: "subscriptionID1",
				hostname:       "hostname1",
			},
			expects: expects{
				GetSubscriptionsReturn: map[string]*types.SubscriptionWithHostname{},
				GetSubscription:        true,
				GetSubscriptionTimes:   2,
				GetSubscriptionRequest: "subscriptionID1",
				GetSubscriptionReturn:  nil,
			},
			want: types.NewNotFoundLookup().Value,
		},
		{
			name: "when dao return error expected error",
			args: args{
				subscriptionID: "subscriptionID1",
				hostname:       "hostname1",
			},
			expects: expects{
				GetSubscriptionsReturn: map[string]*types.SubscriptionWithHostname{},
				GetSubscription:        true,
				GetSubscriptionTimes:   2,
				GetSubscriptionRequest: "subscriptionID1",
				GetSubscriptionError:   errors.New("error1"),
			},
			want:    types.NewNotFoundLookup().Value,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			daoService := daomocks.NewMockIDAO(ctrl)
			s := New(daoService)

			daoService.EXPECT().GetSubscriptions().Return(tt.expects.GetSubscriptionsReturn, tt.expects.GetSubscriptionsError)
			if tt.expects.GetSubscription {
				times := tt.expects.GetSubscriptionTimes
				if times == 0 {
					times = 1
				}
				daoService.EXPECT().GetSubscription(tt.expects.GetSubscriptionRequest).Times(times).Return(tt.expects.GetSubscriptionReturn, tt.expects.GetSubscriptionError)
			}

			err := s.Load()
			if tt.wantGetSubscriptionsErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			var got *types.Lookup
			got, err = s.GetSubscription(tt.args.subscriptionID, tt.args.hostname)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.want, got)

			// should hit the cache the second time
			got, err = s.GetSubscription(tt.args.subscriptionID, tt.args.hostname)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.want, got)
		})
	}
}
