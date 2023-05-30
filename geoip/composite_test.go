package geoip

import (
	"errors"
	"net/http"
	"testing"
)

type testServiceImpl struct {
	country string
	city    string
	err     error
}

// nolint: revive
func (s testServiceImpl) CountryAndCity(_ http.Header, _ string) (string, string, error) {
	return s.country, s.city, s.err
}

func TestComposite_CountryAndCity(t *testing.T) {
	type fields struct {
		primary Service
		next    Service
	}
	type args struct {
		header   http.Header
		ipString string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		{
			name: "primary - country/city",
			args: args{},
			fields: fields{
				primary: testServiceImpl{
					country: "BG",
					city:    "Sofia",
				},
			},
			want:  "BG",
			want1: "Sofia",
		},
		{
			name: "primary - country",
			args: args{},
			fields: fields{
				primary: testServiceImpl{
					country: "BG",
					city:    "",
				},
			},
			want:  "BG",
			want1: "",
		},
		{
			name: "primary - city",
			args: args{},
			fields: fields{
				primary: testServiceImpl{
					country: "",
					city:    "Sofia",
				},
			},
			want:  "",
			want1: "Sofia",
		},
		{
			name: "primary - empty",
			args: args{},
			fields: fields{
				primary: testServiceImpl{
					country: "",
					city:    "",
				},
				next: testServiceImpl{
					country: "BG",
					city:    "Sofia",
				},
			},
			want:  "BG",
			want1: "Sofia",
		},
		{
			name: "primary - error",
			args: args{},
			fields: fields{
				primary: testServiceImpl{
					country: "DE",
					city:    "Berlin",
					err:     errors.New("test"),
				},
				next: testServiceImpl{
					country: "BG",
					city:    "Sofia",
				},
			},
			want:  "BG",
			want1: "Sofia",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Composite{
				primary: tt.fields.primary,
				next:    tt.fields.next,
			}
			got, got1, err := s.CountryAndCity(tt.args.header, tt.args.ipString)
			if (err != nil) != tt.wantErr {
				t.Errorf("Composite.CountryAndCity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Composite.CountryAndCity() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Composite.CountryAndCity() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
