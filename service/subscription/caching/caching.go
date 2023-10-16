package caching

import (
	"github.com/basicrum/front_basicrum_go/dao"
	"github.com/basicrum/front_basicrum_go/service"
	"github.com/basicrum/front_basicrum_go/types"
)

type CachingSubscriptionService struct {
	dao   dao.IDAO
	cache map[string]*types.SubscriptionWithHostname
}

// New creates caching subscription service
func New(daoService dao.IDAO) *CachingSubscriptionService {
	return &CachingSubscriptionService{
		dao:   daoService,
		cache: make(map[string]*types.SubscriptionWithHostname),
	}
}

// Load loads subscriptions from dao to cache
func (s *CachingSubscriptionService) Load() error {
	var err error
	s.cache, err = s.dao.GetSubscriptions()
	return err
}

// GetSubscription attempts to get subscription from cache
// If not successful it attempts to load from dao and updates cache
func (s *CachingSubscriptionService) GetSubscription(subscriptionID, hostname string) (service.Lookup, error) {
	item := s.cache[subscriptionID]
	if item != nil {
		return s.makeLookupResult(item, hostname)
	}

	subscriptionFromDB, err := s.dao.GetSubscription(subscriptionID)
	if err != nil {
		return service.NotFoundLookup, err
	}
	if subscriptionFromDB == nil {
		return service.NotFoundLookup, nil
	}

	s.cache[subscriptionID] = item
	return s.makeLookupResult(item, hostname)
}

func (*CachingSubscriptionService) makeLookupResult(item *types.SubscriptionWithHostname, hostname string) (service.Lookup, error) {
	if item.Subscription.Expired() {
		return service.ExpiredLookup, nil
	}
	if item.Hostname != hostname {
		return service.NotFoundLookup, nil
	}
	return service.FoundLookup, nil
}
