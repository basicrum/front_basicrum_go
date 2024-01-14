package server

import (
	"context"
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

func makeURL(port, path string) string {
	return fmt.Sprintf("http://localhost:%s%s", port, path)
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

func makeServer(processService *servicemocks.MockIService, backupService backup.IBackup) (string, *Server) {
	port := randomPort()
	s := New(processService, backupService, WithHTTP(port))
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
			port, s := makeServer(processService, backupService)

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
