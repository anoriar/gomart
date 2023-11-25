package fetcher

import (
	"github.com/anoriar/gophermart/internal/gophermart/order/entity"
)

type OrderFetchServiceInterface interface {
	Fetch(orderEntity entity.Order) (entity.Order, error)
}
