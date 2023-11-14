package fetcher

import "github.com/anoriar/gophermart/internal/gophermart/entity/order"

type OrderFetchServiceInterface interface {
	Fetch(orderEntity order.Order) (order.Order, error)
}
