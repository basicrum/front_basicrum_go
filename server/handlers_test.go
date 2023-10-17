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
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/basicrum/front_basicrum_go/backup"
	backupmocks "github.com/basicrum/front_basicrum_go/backup/mocks"
	servicemocks "github.com/basicrum/front_basicrum_go/service/mocks"
	"github.com/basicrum/front_basicrum_go/types"
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
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			processService := servicemocks.NewMockIService(ctrl)
			backupService := backup.NewNullBackup()
			port, s := makeServer(processService, backupService, privateAPIToken)

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

func makeUrlValues(beaconDataMap map[string]string) url.Values {
	result := url.Values{}
	for k, v := range beaconDataMap {
		result.Set(k, v)
	}
	return result
}

func makeFormRequest(t *testing.T, address string, pairs map[string]string) *http.Request {
	params := makeUrlValues(pairs)
	req, err := http.NewRequest(http.MethodPost, address, strings.NewReader(params.Encode()))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
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

func makeServer(processService *servicemocks.MockIService, backupService backup.IBackup, privateAPIToken string) (string, *Server) {
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

type eventMatcher struct {
	x *types.Event
}

func (e eventMatcher) Matches(x any) bool {
	arg, ok := x.(*types.Event)
	if !ok {
		return false
	}
	result := reflect.DeepEqual(e.x.RequestParameters, arg.RequestParameters)
	if !result {
		panic(fmt.Sprintf("expected[%v] got[%v]", e.x.RequestParameters, arg.RequestParameters))
	}
	return result
}

func (e eventMatcher) String() string {
	return fmt.Sprintf("is equal to %v", e.x)
}

func eqEvent(arg *types.Event) gomock.Matcher {
	return eventMatcher{arg}
}

func TestServer_catcher(t *testing.T) {
	requestForm := map[string]string{
		"hostname":        "hostname1",
		"subscription_id": "subscription_id1",
		"created_at":      "created_at1",
	}
	expectedEvent := &types.Event{
		RequestParameters: url.Values{
			"hostname":        []string{"hostname1"},
			"subscription_id": []string{"subscription_id1"},
			"created_at":      []string{"created_at1"},
		},
		UserAgent: "Go-http-client/1.1",
	}
	type args struct {
		form map[string]string
	}
	type expects struct {
		SaveAsync              bool
		SaveAsyncRequest       *types.Event
		BackupSaveAsync        bool
		BackupSaveAsyncRequest *types.Event
	}
	tests := []struct {
		name     string
		args     args
		expects  expects
		want     string
		wantCode int
	}{
		{
			name: "Success",
			args: args{
				form: requestForm,
			},
			expects: expects{
				SaveAsync:              true,
				SaveAsyncRequest:       expectedEvent,
				BackupSaveAsync:        true,
				BackupSaveAsyncRequest: expectedEvent,
			},
			want:     "",
			wantCode: http.StatusNoContent,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			processService := servicemocks.NewMockIService(ctrl)
			backupService := backupmocks.NewMockIBackup(ctrl)
			port, s := makeServer(processService, backupService, "")

			go func() {
				_ = s.Serve()
			}()
			defer func() {
				_ = s.Shutdown(context.Background())
			}()
			if tt.expects.SaveAsync {
				processService.EXPECT().SaveAsync(eqEvent(tt.expects.SaveAsyncRequest))
			}
			if tt.expects.BackupSaveAsync {
				backupService.EXPECT().SaveAsync(eqEvent(tt.expects.BackupSaveAsyncRequest))
			}
			address := makeURL(port, "/beacon/catcher")
			r := makeFormRequest(t, address, tt.args.form)
			response := executeRequest(r, t)

			assertResponse(t, response, tt.want, tt.wantCode)
		})
	}
}
