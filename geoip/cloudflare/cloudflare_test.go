package cloudflare

import (
	"net/http"
	"testing"
)

func TestService_CountryAndCity(t *testing.T) {
	header1 := http.Header{}
	header1.Set("CF-IPCountry", "BG")
	header1.Set("CF-IPCity", "Sofia")

	type args struct {
		header http.Header
	}
	tests := []struct {
		name    string
		s       *Service
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				header: header1,
			},
			want:  "BG",
			want1: "Sofia",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New()
			got, got1, err := s.CountryAndCity(tt.args.header, "")
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.CountryAndCity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Service.CountryAndCity() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Service.CountryAndCity() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
