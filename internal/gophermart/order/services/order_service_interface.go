package services

import (
	"context"
	orderRepositoryDtoPkg "github.com/anoriar/gophermart/internal/gophermart/order/dto/repository"
	"github.com/anoriar/gophermart/internal/gophermart/order/entity"
)

type OrderServiceInterface interface {
	LoadOrder(ctx context.Context, orderID string, userID string) error
	GetUserOrders(ctx context.Context, userID string) ([]entity.Order, error)
	ProcessOrder(ctx context.Context, orderID string) error
	GetOrders(ctx context.Context, query orderRepositoryDtoPkg.OrdersQuery) ([]entity.Order, error)
	GetOrdersTotal(ctx context.Context, filter orderRepositoryDtoPkg.OrdersFilterDto) (int, error)
}
