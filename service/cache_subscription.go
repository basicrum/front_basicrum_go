package service

import (
	"time"

	"github.com/basicrum/front_basicrum_go/dao"
	"github.com/basicrum/front_basicrum_go/types"
)

type CacheSubscriptionService struct {
	dao   dao.DAO
	cache map[string]types.Subscription
}

// NewCacheSubscriptionService creates caching subscription service
func NewCacheSubscriptionService(dao dao.DAO) *CacheSubscriptionService {
	return &CacheSubscriptionService{
		dao:   dao,
		cache: make(map[string]types.Subscription),
	}
}

// Load loads subscriptions from dao to cache
func (s *CacheSubscriptionService) Load() error {
	var err error
	s.cache, err = s.dao.GetSubscriptions()
	if err != nil {
		return err
	}
	return nil
}

// Update updates cache by key
func (s *CacheSubscriptionService) Update(key string, subscription types.Subscription) {
	s.cache[key] = subscription
}

// GetSubscription attempts to get subscription from cache
// If not successful it attempts to load from dao and updates cache
func (s *CacheSubscriptionService) GetSubscription(id string) (Lookup, error) {
	if subscription, ok := s.cache[id]; ok {
		if subscriptionExpired(subscription.ExpiresAt) {
			return ExpiredLookup, nil
		}
		return FoundLookup, nil
	}

	subscription, err := s.dao.GetSubscription(id)
	if err != nil {
		return NotFoundLookup, err
	}

	s.Update(subscription.ID, subscription)
	if subscriptionExpired(subscription.ExpiresAt) {
		return ExpiredLookup, nil
	} else {
		return FoundLookup, nil
	}
}

func subscriptionExpired(expiresAt time.Time) bool {
	return time.Now().Before(expiresAt)
}
