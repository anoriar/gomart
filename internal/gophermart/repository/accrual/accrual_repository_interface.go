package accrual

import "github.com/anoriar/gophermart/internal/gophermart/dto/accrual"

type AccrualRepositoryInterface interface {
	GetOrder(orderId string) (result accrual.AccrualOrderDto, exists bool, err error)
}
