package service

import (
	"github.com/basicrum/front_basicrum_go/dao"
	"github.com/basicrum/front_basicrum_go/types"
)

type CacheSubscriptionService struct {
	dao   dao.DAO
	cache map[string]types.Subscription
}

func NewCacheSubscriptionService(dao dao.DAO) *CacheSubscriptionService {
	return &CacheSubscriptionService{
		dao:   dao,
		cache: make(map[string]types.Subscription),
	}
}

func (s *CacheSubscriptionService) Load() error {
	// TODO: call dao
	return nil
}

func (s *CacheSubscriptionService) GetSubscription(id string) (Lookup, error) {
	// TODO: implement
	// search in cache
	// if found in cache return
	// load from dao
	// update cache
	return NotFoundLookup, nil
}
