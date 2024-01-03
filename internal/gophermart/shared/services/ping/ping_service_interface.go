package ping

import (
	"context"
	"github.com/anoriar/gophermart/internal/gophermart/shared/dto/responses/ping"
)

//go:generate mockgen -source=ping_service_interface.go -destination=mock/ping_service.go -package=mock
type PingServiceInterface interface {
	Ping(ctx context.Context) ping.PingResponseDto
}
