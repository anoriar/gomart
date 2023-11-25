package repository

import (
	"context"
	orderRepositoryDtoPkg "github.com/anoriar/gophermart/internal/gophermart/order/dto/repository"
	"github.com/anoriar/gophermart/internal/gophermart/order/entity"
)

type OrderRepositoryInterface interface {
	AddOrder(ctx context.Context, order entity.Order) error
	GetOrderByID(ctx context.Context, orderID string) (entity.Order, error)
	GetOrders(ctx context.Context, query orderRepositoryDtoPkg.OrdersQuery) ([]entity.Order, error)
	GetTotal(ctx context.Context, filter orderRepositoryDtoPkg.OrdersFilterDto) (int, error)
	UpdateOrder(ctx context.Context, orderID string, status string, accrual float64) error
}
