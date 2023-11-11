package order

import (
	"context"
	"github.com/anoriar/gophermart/internal/gophermart/entity/order"
)

type OrderRepositoryInterface interface {
	AddOrder(ctx context.Context, order order.Order) error
	GetOrderByID(ctx context.Context, orderID string) (order.Order, error)
}
