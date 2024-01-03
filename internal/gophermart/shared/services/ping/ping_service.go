package ping

import (
	"context"
	db2 "github.com/anoriar/gophermart/internal/gophermart/shared/app/db"
	ping2 "github.com/anoriar/gophermart/internal/gophermart/shared/dto/responses/ping"
	"github.com/opentracing/opentracing-go"
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

func (service *PingService) Ping(ctx context.Context) ping2.PingResponseDto {
	return ping2.PingResponseDto{
		Services: []ping2.ServiceStatusDto{
			service.pingDatabase(ctx),
		},
	}
}

func (service *PingService) pingDatabase(ctx context.Context) ping2.ServiceStatusDto {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PingService::pingDatabase")
	defer span.Finish()

	err := service.database.Ping(ctx)
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
