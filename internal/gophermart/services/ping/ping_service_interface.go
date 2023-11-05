package ping

import "github.com/anoriar/gophermart/internal/gophermart/dto/responses/ping"

//go:generate mockgen -source=ping_service_interface.go -destination=mock/ping_service.go -package=mock
type PingServiceInterface interface {
	Ping() ping.PingResponseDto
}
