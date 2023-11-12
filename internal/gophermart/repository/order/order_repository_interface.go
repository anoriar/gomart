package order

import (
	"context"
	orderQueryPkg "github.com/anoriar/gophermart/internal/gophermart/dto/repository/order"
	"github.com/anoriar/gophermart/internal/gophermart/entity/order"
)

type OrderRepositoryInterface interface {
	AddOrder(ctx context.Context, order order.Order) error
	GetOrderByID(ctx context.Context, orderID string) (order.Order, error)
	GetOrders(ctx context.Context, query orderQueryPkg.OrdersQuery) ([]order.Order, error)
	GetTotal(ctx context.Context, filter orderQueryPkg.OrdersFilterDto) (int, error)
}
