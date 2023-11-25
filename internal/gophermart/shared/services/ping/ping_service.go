package ping

import (
	db2 "github.com/anoriar/gophermart/internal/gophermart/shared/app/db"
	ping2 "github.com/anoriar/gophermart/internal/gophermart/shared/dto/responses/ping"
)

const (
	dbServiceName = "Database"
)

type PingService struct {
	database db2.DatabaseInterface
}

func NewPingService(database *db2.Database) *PingService {
	return &PingService{database: database}
}

func (service *PingService) Ping() ping2.PingResponseDto {
	return ping2.PingResponseDto{
		Services: []ping2.ServiceStatusDto{
			service.pingDatabase(),
		},
	}
}

func (service *PingService) pingDatabase() ping2.ServiceStatusDto {
	err := service.database.Ping()
	if err != nil {
		return ping2.ServiceStatusDto{
			Name:   dbServiceName,
			Status: ping2.FailStatus,
			Error:  err.Error(),
		}
	}
	return ping2.ServiceStatusDto{
		Name:   dbServiceName,
		Status: ping2.OKStatus,
		Error:  "",
	}
}
