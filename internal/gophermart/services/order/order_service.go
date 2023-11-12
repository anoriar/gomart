package order

import (
	"context"
	"errors"
	"fmt"
	"github.com/anoriar/gophermart/internal/gophermart/domain_errors"
	orderQueryPkg "github.com/anoriar/gophermart/internal/gophermart/dto/repository/order"
	orderPkg "github.com/anoriar/gophermart/internal/gophermart/entity/order"
	"github.com/anoriar/gophermart/internal/gophermart/repository/order"
	"github.com/anoriar/gophermart/internal/gophermart/services/order/internal/services/luhn_validator"
	"go.uber.org/zap"
)

type OrderService struct {
	orderRepository order.OrderRepositoryInterface
	luhnValidator   *luhn_validator.LuhnValidator
	logger          *zap.Logger
}

func (service *OrderService) GetUserOrders(ctx context.Context, userID string) ([]orderPkg.Order, error) {
	orders, err := service.orderRepository.GetOrders(ctx, orderQueryPkg.OrdersQuery{
		Filter: orderQueryPkg.OrdersFilterDto{
			UserID: userID,
		},
	})
	if err != nil {
		service.logger.Error(err.Error())
		return nil, err
	}

	return orders, nil
}

func NewOrderService(orderRepository order.OrderRepositoryInterface, logger *zap.Logger) *OrderService {
	return &OrderService{orderRepository: orderRepository, luhnValidator: luhn_validator.NewLuhnValidator(), logger: logger}
}

func (service *OrderService) LoadOrder(ctx context.Context, orderID string, userID string) error {
	if service.luhnValidator.Validate(orderID) == false {
		return domain_errors.ErrOrderNumberNotValid
	}
	currentOrder, err := service.orderRepository.GetOrderByID(ctx, orderID)

	switch {
	case err != nil && errors.Is(err, domain_errors.ErrNotFound):
		newOrder := orderPkg.CreateNewOrder(orderID, userID)
		err := service.orderRepository.AddOrder(ctx, newOrder)
		if err != nil {
			service.logger.Error(err.Error())
			return err
		}
		//TODO: send async task
		return nil
	case err != nil && !errors.Is(err, domain_errors.ErrNotFound):
		service.logger.Error(err.Error())
		return err
	case currentOrder.UserID == userID:
		return domain_errors.ErrOrderAlreadyLoaded
	case currentOrder.UserID != userID:
		return domain_errors.ErrOrderLoadedByAnotherUser
	default:
		service.logger.Error("unexpected behaviour")
		return fmt.Errorf("unexpected behaviour %w", domain_errors.ErrInternalError)
	}
}
