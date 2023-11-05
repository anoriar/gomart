package ping

import (
	"github.com/anoriar/gophermart/internal/gophermart/app/db"
	"github.com/anoriar/gophermart/internal/gophermart/dto/responses/ping"
)

const (
	dbServiceName = "Database"
)

type PingService struct {
	database db.DatabaseInterface
}

func NewPingService(database *db.Database) *PingService {
	return &PingService{database: database}
}

func (service *PingService) Ping() ping.PingResponseDto {
	return ping.PingResponseDto{
		Services: []ping.ServiceStatus{
			service.pingDatabase(),
		},
	}
}

func (service *PingService) pingDatabase() ping.ServiceStatus {
	err := service.database.Ping()
	if err != nil {
		return ping.ServiceStatus{
			Name:   dbServiceName,
			Status: ping.FailStatus,
			Error:  err.Error(),
		}
	}
	return ping.ServiceStatus{
		Name:   dbServiceName,
		Status: ping.OKStatus,
		Error:  "",
	}
}
