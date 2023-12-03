package ping

import (
	"errors"
	"github.com/anoriar/gophermart/internal/gophermart/shared/app/db/mock"
	ping2 "github.com/anoriar/gophermart/internal/gophermart/shared/dto/responses/ping"
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
		want          ping2.PingResponseDto
	}{
		{
			name: "success",
			mockBehaviour: func() {
				dbMock.EXPECT().Ping().Return(nil).Times(1)
			},
			want: ping2.PingResponseDto{
				Services: []ping2.ServiceStatusDto{
					{
						Name:   dbServiceName,
						Status: ping2.OKStatus,
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
			want: ping2.PingResponseDto{
				Services: []ping2.ServiceStatusDto{
					{
						Name:   dbServiceName,
						Status: ping2.FailStatus,
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
