package factory

import (
	orderRsponseDtoPkg "github.com/anoriar/gophermart/internal/gophermart/dto/responses/order"
	"github.com/anoriar/gophermart/internal/gophermart/entity/order"
	"time"
)

type GetOrdersResponseFactory struct {
}

func NewGetOrdersResponseFactory() *GetOrdersResponseFactory {
	return &GetOrdersResponseFactory{}
}

func (factory GetOrdersResponseFactory) CreateOrdersResponse(orders []order.Order) []orderRsponseDtoPkg.OrderResponseDto {
	var response []orderRsponseDtoPkg.OrderResponseDto
	for _, orderEntity := range orders {
		response = append(response, orderRsponseDtoPkg.OrderResponseDto{
			Number:     orderEntity.ID,
			Status:     orderEntity.Status,
			Accrual:    orderEntity.Accrual,
			UploadedAt: orderEntity.UploadedAt.Format(time.RFC3339),
		})
	}
	return response
}
