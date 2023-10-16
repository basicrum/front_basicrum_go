package caching

import (
	"testing"

	"github.com/basicrum/front_basicrum_go/dao"
	"github.com/basicrum/front_basicrum_go/service"
	"github.com/basicrum/front_basicrum_go/types"
)

func TestCacheSubscriptionService_GetSubscription(t *testing.T) {
	type fields struct {
		dao   dao.IDAO
		cache map[string]*types.SubscriptionWithHostname
	}
	type args struct {
		subscriptionID string
		hostname       string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    service.Lookup
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &CachingSubscriptionService{
				dao:   tt.fields.dao,
				cache: tt.fields.cache,
			}
			got, err := s.GetSubscription(tt.args.subscriptionID, tt.args.hostname)
			if (err != nil) != tt.wantErr {
				t.Errorf("CacheSubscriptionService.GetSubscription() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CacheSubscriptionService.GetSubscription() = %v, want %v", got, tt.want)
			}
		})
	}
}
