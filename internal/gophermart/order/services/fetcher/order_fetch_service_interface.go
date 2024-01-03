package fetcher

import (
	"context"
	"github.com/anoriar/gophermart/internal/gophermart/order/entity"
)

type OrderFetchServiceInterface interface {
	Fetch(ctx context.Context, orderEntity entity.Order) (entity.Order, error)
}
