package services

import (
	"context"
	"github.com/anoriar/gophermart/internal/gophermart/order/entity"
)

type OrderServiceInterface interface {
	LoadOrder(ctx context.Context, orderID string, userID string) error
	GetUserOrders(ctx context.Context, userID string) ([]entity.Order, error)
	ProcessOrder(ctx context.Context, orderID string) error
}
