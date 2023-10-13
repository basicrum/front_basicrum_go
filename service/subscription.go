package service

type Lookup string

const (
	FoundLookup    Lookup = "FOUND"
	NotFoundLookup Lookup = "NOT_FOUND"
	ExpiredLookup  Lookup = "EXPIRED"
)

type ISubscriptionService interface {
	Load() error
	GetSubscription(id string) (Lookup, error)
}
