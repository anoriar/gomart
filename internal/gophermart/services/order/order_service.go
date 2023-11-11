package order

import (
	"context"
	"errors"
	"fmt"
	"github.com/anoriar/gophermart/internal/gophermart/domain_errors"
	order2 "github.com/anoriar/gophermart/internal/gophermart/entity/order"
	"github.com/anoriar/gophermart/internal/gophermart/repository/order"
	"github.com/anoriar/gophermart/internal/gophermart/services/order/internal/services/luhn_validator"
	"go.uber.org/zap"
)

type OrderService struct {
	orderRepository order.OrderRepositoryInterface
	luhnValidator   *luhn_validator.LuhnValidator
	logger          *zap.Logger
}

func NewOrderService(orderRepository order.OrderRepositoryInterface, luhnValidator *luhn_validator.LuhnValidator, logger *zap.Logger) *OrderService {
	return &OrderService{orderRepository: orderRepository, luhnValidator: luhnValidator, logger: logger}
}

func (service *OrderService) LoadOrder(ctx context.Context, orderID string, userID string) error {
	if service.luhnValidator.Validate(orderID) == false {
		return domain_errors.ErrOrderNumberNotValid
	}
	currentOrder, err := service.orderRepository.GetOrderByID(ctx, orderID)

	switch {
	case err != nil && errors.Is(err, domain_errors.ErrNotFound):
		newOrder := order2.CreateNewOrder(orderID, userID)
		err := service.orderRepository.AddOrder(ctx, newOrder)
		if err != nil {
			service.logger.Error(err.Error())
			return err
		}
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
