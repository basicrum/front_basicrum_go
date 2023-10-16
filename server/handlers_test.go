package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/basicrum/front_basicrum_go/backup"
	servicemocks "github.com/basicrum/front_basicrum_go/service/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestServer_hostnames(t *testing.T) {
	privateAPIToken := "privateAPIToken1"
	type args struct {
		method  string
		request map[string]any
		headers map[string]string
	}
	type expects struct {
		RegisterHostname bool
		DeleteHostname   bool
		hostname         string
		username         string
		ReturnError      error
	}
	tests := []struct {
		name     string
		args     args
		expects  expects
		want     string
		wantCode int
	}{
		{
			name: "POST success",
			args: args{
				method: http.MethodPost,
				request: map[string]any{
					"hostname": "test1",
				},
				headers: map[string]string{
					"X-Grafana-User": "user1",
					"X-Token":        privateAPIToken,
				},
			},
			expects: expects{
				RegisterHostname: true,
				hostname:         "test1",
				username:         "user1",
			},
			want:     "ok",
			wantCode: http.StatusOK,
		},
		{
			name: "POST error",
			args: args{
				method: http.MethodPost,
				request: map[string]any{
					"hostname": "test1",
				},
				headers: map[string]string{
					"X-Grafana-User": "user1",
					"X-Token":        privateAPIToken,
				},
			},
			expects: expects{
				RegisterHostname: true,
				hostname:         "test1",
				username:         "user1",
				ReturnError:      errors.New("error1"),
			},
			want:     "error1",
			wantCode: http.StatusBadRequest,
		},
		{
			name: "POST hostname is required",
			args: args{
				method: http.MethodPost,
				request: map[string]any{
					"hostname": "",
				},
				headers: map[string]string{
					"X-Grafana-User": "user1",
					"X-Token":        privateAPIToken,
				},
			},
			want:     "hostname is required",
			wantCode: http.StatusBadRequest,
		},
		{
			name: "POST required X-Grafana-User",
			args: args{
				method: http.MethodPost,
				request: map[string]any{
					"hostname": "test1",
				},
				headers: map[string]string{
					"X-Grafana-User": "",
					"X-Token":        privateAPIToken,
				},
			},
			want:     "required header[X-Grafana-User]",
			wantCode: http.StatusBadRequest,
		},
		{
			name: "POST WRONG_TOKEN",
			args: args{
				method: http.MethodPost,
				headers: map[string]string{
					"X-Token": "WRONG_TOKEN",
				},
			},
			want:     "wrong header[X-Token]",
			wantCode: http.StatusBadRequest,
		},
		{
			name: "DELETE success",
			args: args{
				method: http.MethodDelete,
				request: map[string]any{
					"hostname": "test1",
				},
				headers: map[string]string{
					"X-Grafana-User": "user1",
					"X-Token":        privateAPIToken,
				},
			},
			expects: expects{
				DeleteHostname: true,
				hostname:       "test1",
				username:       "user1",
			},
			want:     "ok",
			wantCode: http.StatusOK,
		},
		{
			name: "DELETE error",
			args: args{
				method: http.MethodDelete,
				request: map[string]any{
					"hostname": "test1",
				},
				headers: map[string]string{
					"X-Grafana-User": "user1",
					"X-Token":        privateAPIToken,
				},
			},
			expects: expects{
				DeleteHostname: true,
				hostname:       "test1",
				username:       "user1",
				ReturnError:    errors.New("error1"),
			},
			want:     "error1",
			wantCode: http.StatusBadRequest,
		},
		{
			name: "DELETE hostname is required",
			args: args{
				method: http.MethodDelete,
				request: map[string]any{
					"hostname": "",
				},
				headers: map[string]string{
					"X-Grafana-User": "user1",
					"X-Token":        privateAPIToken,
				},
			},
			want:     "hostname is required",
			wantCode: http.StatusBadRequest,
		},
		{
			name: "DELETE required X-Grafana-User",
			args: args{
				method: http.MethodDelete,
				request: map[string]any{
					"hostname": "test1",
				},
				headers: map[string]string{
					"X-Grafana-User": "",
					"X-Token":        privateAPIToken,
				},
			},
			want:     "required header[X-Grafana-User]",
			wantCode: http.StatusBadRequest,
		},
		{
			name: "DELETE WRONG_TOKEN",
			args: args{
				method: http.MethodDelete,
				headers: map[string]string{
					"X-Token": "WRONG_TOKEN",
				},
			},
			want:     "wrong header[X-Token]",
			wantCode: http.StatusBadRequest,
		},
		{
			name: "PUT unsupported method",
			args: args{
				method: http.MethodPut,
				headers: map[string]string{
					"X-Token": privateAPIToken,
				},
			},
			want:     "unsupported method[PUT]",
			wantCode: http.StatusBadRequest,
		},
		{
			name: "GET unsupported method",
			args: args{
				method: http.MethodGet,
				headers: map[string]string{
					"X-Token": privateAPIToken,
				},
			},
			want:     "unsupported method[GET]",
			wantCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// add a sleep of one second between spawning servers to avoid connection refused on slower cpus
			time.Sleep(time.Second)
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			processService := servicemocks.NewMockIService(ctrl)
			port, s := makeServer(processService, privateAPIToken)

			go func() {
				_ = s.Serve()
			}()
			defer func() {
				_ = s.Shutdown(context.Background())
			}()

			if tt.expects.RegisterHostname {
				processService.EXPECT().
					RegisterHostname(
						tt.expects.hostname,
						tt.expects.username,
					).Return(tt.expects.ReturnError)
			}
			if tt.expects.DeleteHostname {
				processService.EXPECT().
					DeleteHostname(
						tt.expects.hostname,
						tt.expects.username,
					).Return(tt.expects.ReturnError)
			}
			address := makeURL(port, "/private/hostnames")
			r := makeRequest(t, tt.args.method, address, tt.args.request, tt.args.headers)

			response := executeRequest(r, t)

			assertResponse(t, response, tt.want, tt.wantCode)
		})
	}
}

func makeURL(port, path string) string {
	return fmt.Sprintf("http://localhost:%s%s", port, path)
}

// nolint: revive
func makeRequest(t *testing.T, method, address string, request any, headers map[string]string) *http.Request {
	b, err := json.Marshal(request)
	require.NoError(t, err)

	r, err := http.NewRequest(method, address, bytes.NewBuffer(b))
	require.NoError(t, err)

	r.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		r.Header.Set(k, v)
	}
	return r
}

func executeRequest(r *http.Request, t *testing.T) *http.Response {
	response, err := http.DefaultClient.Do(r)
	require.NoError(t, err)
	return response
}

func assertResponse(t *testing.T, response *http.Response, want string, wantCode int) {
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	require.Equal(t, want, string(responseBody))
	require.Equal(t, wantCode, response.StatusCode)
}

func makeServer(processService *servicemocks.MockIService, privateAPIToken string) (string, *Server) {
	backupService := backup.NewNullBackup()
	port := randomPort()
	s := New(processService, backupService, privateAPIToken, WithHTTP(port))
	return port, s
}

func randomPort() string {
	srv := httptest.NewServer(nil)
	defer srv.Close()
	pos := strings.LastIndex(srv.URL, ":")
	return srv.URL[pos+1:]
}

// func TestServer_catcher(t *testing.T) {
// 	type fields struct {
// 		port            string
// 		service         service.IService
// 		backup          backup.IBackup
// 		certFile        string
// 		keyFile         string
// 		server          *http.Server
// 		tlsConfig       *tls.Config
// 		privateAPIToken string
// 	}
// 	type args struct {
// 		w http.ResponseWriter
// 		r *http.Request
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := &Server{
// 				port:            tt.fields.port,
// 				service:         tt.fields.service,
// 				backup:          tt.fields.backup,
// 				certFile:        tt.fields.certFile,
// 				keyFile:         tt.fields.keyFile,
// 				server:          tt.fields.server,
// 				tlsConfig:       tt.fields.tlsConfig,
// 				privateAPIToken: tt.fields.privateAPIToken,
// 			}
// 			s.catcher(tt.args.w, tt.args.r)
// 		})
// 	}
// }
