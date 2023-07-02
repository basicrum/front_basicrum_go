package beacon

import (
	"net/http"
	"testing"

	"github.com/basicrum/front_basicrum_go/types"
	"github.com/stretchr/testify/assert"
	"github.com/ua-parser/uap-go/uaparser"
)

type mockGeoIPService struct{}

// nolint: revive
func (*mockGeoIPService) CountryAndCity(_ http.Header, _ string) (string, string, error) {
	return "", "", nil
}

type mockUserAgentParser struct{}

func (*mockUserAgentParser) Parse(_ string) *uaparser.Client {
	return &uaparser.Client{
		UserAgent: &uaparser.UserAgent{},
		Os:        &uaparser.Os{},
		Device:    &uaparser.Device{},
	}
}

func TestConvertToRumEvent(t *testing.T) {
	type args struct {
		b     Beacon
		event *types.Event
	}
	tests := []struct {
		name string
		args args
		want RumEvent
	}{
		{
			name: "default",
			args: args{
				b:     Beacon{},
				event: &types.Event{},
			},
			want: RumEvent{
				Device_Type:       "unknown",
				Event_Type:        "visit_page",
				Redirect_Duration: "0",
				Redirects_Count:   "0",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertToRumEvent(tt.args.b, tt.args.event, &mockUserAgentParser{}, &mockGeoIPService{})
			assert.Equal(t, tt.want, got)
		})
	}
}
