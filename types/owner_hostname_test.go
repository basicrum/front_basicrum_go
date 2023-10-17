package types

import (
	"testing"
	"time"
)

func TestSubscription_Expired(t *testing.T) {
	type fields struct {
		ID        string
		ExpiresAt time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "not expired - valid",
			fields: fields{
				ExpiresAt: time.Now().Add(time.Hour),
			},
			want: false,
		},
		{
			name: "expired",
			fields: fields{
				ExpiresAt: time.Now().Add(-time.Hour),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Subscription{
				ID:        tt.fields.ID,
				ExpiresAt: tt.fields.ExpiresAt,
			}
			if got := s.Expired(); got != tt.want {
				t.Errorf("Subscription.Expired() = %v, want %v", got, tt.want)
			}
		})
	}
}
