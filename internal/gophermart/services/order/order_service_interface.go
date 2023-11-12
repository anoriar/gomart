package order

import (
	"context"
	"github.com/anoriar/gophermart/internal/gophermart/entity/order"
)

type OrderServiceInterface interface {
	LoadOrder(ctx context.Context, orderID string, userID string) error
	GetUserOrders(ctx context.Context, userID string) ([]order.Order, error)
}
