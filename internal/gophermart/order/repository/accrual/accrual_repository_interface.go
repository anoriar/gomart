package accrual

import (
	"context"
	"github.com/anoriar/gophermart/internal/gophermart/order/dto/accrual"
)

type AccrualRepositoryInterface interface {
	GetOrder(ctx context.Context, orderID string) (result accrual.AccrualOrderDto, exists bool, err error)
}
