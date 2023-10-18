package service

import "github.com/basicrum/front_basicrum_go/types"

//go:generate mockgen -source=${GOFILE} -destination=mocks/${GOFILE} -package=servicemocks

// Lookup describes subscription lookup statuses
// type Lookup string

// const (
// 	// FoundLookup found
// 	FoundLookup Lookup = "FOUND"
// 	// NotFoundLookup not found
// 	NotFoundLookup Lookup = "NOT_FOUND"
// 	// ExpiredLookup expired
// 	ExpiredLookup Lookup = "EXPIRED"
// )

// ISubscriptionService subscription service
type ISubscriptionService interface {
	// Load initial data
	Load() error
	// GetSubscription get subscription by id and hostname
	GetSubscription(subscriptionID, hostname string) (*types.Lookup, error)
}
