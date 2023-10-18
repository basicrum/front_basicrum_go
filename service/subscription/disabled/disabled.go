package disabled

import (
	"github.com/basicrum/front_basicrum_go/service"
	"github.com/basicrum/front_basicrum_go/types"
)

type disabledSubscriptionService struct {
}

// New creates disabled subscription service
func New() service.ISubscriptionService {
	return &disabledSubscriptionService{}
}

// Load loads subscriptions from dao to cache
func (*disabledSubscriptionService) Load() error {
	return nil
}

// GetSubscription attempts to get subscription from cache
// If not successful it attempts to load from dao and updates cache
func (*disabledSubscriptionService) GetSubscription(_, _ string) (*types.Lookup, error) {
	return types.NewFoundLookup().Value, nil
}
