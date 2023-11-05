package ping

import (
	"errors"
	"github.com/anoriar/gophermart/internal/gophermart/app/db/mock"
	"github.com/anoriar/gophermart/internal/gophermart/dto/responses/ping"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPingService_Ping(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dbError := errors.New("database error")

	dbMock := mock.NewMockDatabaseInterface(ctrl)

	tests := []struct {
		name          string
		mockBehaviour func()
		want          ping.PingResponseDto
	}{
		{
			name: "success",
			mockBehaviour: func() {
				dbMock.EXPECT().Ping().Return(nil).Times(1)
			},
			want: ping.PingResponseDto{
				Services: []ping.ServiceStatus{
					{
						Name:   dbServiceName,
						Status: ping.OKStatus,
						Error:  "",
					},
				},
			},
		},
		{
			name: "database fail",
			mockBehaviour: func() {
				dbMock.EXPECT().Ping().Return(dbError).Times(1)
			},
			want: ping.PingResponseDto{
				Services: []ping.ServiceStatus{
					{
						Name:   dbServiceName,
						Status: ping.FailStatus,
						Error:  dbError.Error(),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour()
			service := &PingService{
				database: dbMock,
			}
			got := service.Ping()
			assert.Equal(t, got, tt.want)
		})
	}
}
