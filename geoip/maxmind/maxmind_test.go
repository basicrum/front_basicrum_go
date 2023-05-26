package maxmind

import (
	"net/http"
	"testing"
)

func TestService_CountryAndCity(t *testing.T) {
	type args struct {
		ipString string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		wantErr bool
	}{

		{
			name: "Country GB without City",
			args: args{
				ipString: "81.2.69.142",
			},
			want:  "GB",
			want1: "",
		},
		{
			name: "City Sofia Bulgaria",
			args: args{
				ipString: "212.5.142.168",
			},
			want:  "BG",
			want1: "Sofia",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New()
			header := http.Header{}
			got, got1, err := s.CountryAndCity(header, tt.args.ipString)
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
