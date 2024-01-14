package types

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	subscriptionTrialMonths = 3
	yearMonthDayFormat      = "2006-01-02"
)

// Subscription represents when the trial for hostname is active
type Subscription struct {
	ID        string
	ExpiresAt time.Time
}

// NewSubscription creates Subscription
func NewSubscription(now time.Time) Subscription {
	expiresAt := now.AddDate(0, subscriptionTrialMonths, 0)
	return Subscription{
		ID:        generateSubscriptionID(expiresAt),
		ExpiresAt: expiresAt,
	}
}

func generateSubscriptionID(expiresAt time.Time) string {
	return fmt.Sprintf("%s|%s", expiresAt.Format(yearMonthDayFormat), uuid.NewString())
}

// OwnerHostname is the hostname registered by owner
type OwnerHostname struct {
	Username     string
	Hostname     string
	Subscription Subscription
}

// NewOwnerHostname creates OwnerHostname
func NewOwnerHostname(
	username string,
	hostname string,
	subscription Subscription,
) OwnerHostname {
	return OwnerHostname{
		Username:     username,
		Hostname:     hostname,
		Subscription: subscription,
	}
}
