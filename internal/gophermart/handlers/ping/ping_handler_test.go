package ping

import (
	"encoding/json"
	"github.com/anoriar/gophermart/internal/gophermart/dto/responses/ping"
	"github.com/anoriar/gophermart/internal/gophermart/services/ping/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingHandler_Ping(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pingServiceMock := mock.NewMockPingServiceInterface(ctrl)
	pingResponseDto := ping.PingResponseDto{
		Services: []ping.ServiceStatusDto{
			{
				Name:   "service",
				Status: ping.OKStatus,
				Error:  "",
			},
		},
	}
	pingResponseBody, err := json.Marshal(pingResponseDto)
	require.NoError(t, err)

	type want struct {
		status      int
		body        []byte
		contentType string
	}
	tests := []struct {
		name          string
		mockBehaviour func()
		want          want
	}{
		{
			name: "success",
			mockBehaviour: func() {
				pingServiceMock.EXPECT().Ping().Return(pingResponseDto)
			},
			want: want{
				status:      http.StatusOK,
				body:        pingResponseBody,
				contentType: "application/json",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour()

			r := httptest.NewRequest(http.MethodGet, "/ping", nil)
			w := httptest.NewRecorder()

			handler := &PingHandler{
				pingService: pingServiceMock,
			}
			handler.Ping(w, r)

			assert.Equal(t, tt.want.status, w.Code)
			assert.Equal(t, tt.want.contentType, w.Header().Get("Content-Type"))

			assert.Equal(t, tt.want.body, w.Body.Bytes())
		})
	}
}
