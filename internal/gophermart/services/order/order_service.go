package order

import (
	"context"
	"errors"
	"fmt"
	"github.com/anoriar/gophermart/internal/gophermart/domain_errors"
	accrualPkg "github.com/anoriar/gophermart/internal/gophermart/dto/accrual"
	"github.com/anoriar/gophermart/internal/gophermart/dto/repository"
	orderQueryPkg "github.com/anoriar/gophermart/internal/gophermart/dto/repository/order"
	orderPkg "github.com/anoriar/gophermart/internal/gophermart/entity/order"
	"github.com/anoriar/gophermart/internal/gophermart/repository/accrual"
	"github.com/anoriar/gophermart/internal/gophermart/repository/order"
	"github.com/anoriar/gophermart/internal/gophermart/services/order/internal/services/luhn_validator"
	"go.uber.org/zap"
)

type OrderService struct {
	orderRepository   order.OrderRepositoryInterface
	accrualRepository accrual.AccrualRepositoryInterface
	luhnValidator     *luhn_validator.LuhnValidator
	logger            *zap.Logger
}

func NewOrderService(orderRepository order.OrderRepositoryInterface, accrualRepository accrual.AccrualRepositoryInterface, logger *zap.Logger) *OrderService {
	return &OrderService{orderRepository: orderRepository, accrualRepository: accrualRepository, luhnValidator: luhn_validator.NewLuhnValidator(), logger: logger}
}

func (service *OrderService) ProcessOrder(ctx context.Context, orderID string) error {
	orderEntity, err := service.orderRepository.GetOrderByID(ctx, orderID)
	if err != nil {
		service.logger.Error(err.Error())
		return err
	}
	if !(orderEntity.Status == orderPkg.NewStatus || orderEntity.Status == orderPkg.ProcessingStatus) {
		errText := fmt.Sprintf("invalid status of order. id: %s", orderID)
		service.logger.Error(errText)
		return fmt.Errorf(errText)
	}

	orderFromAccrualSystem, err := service.accrualRepository.GetOrder(orderID)
	if err != nil {
		if errors.Is(err, domain_errors.ErrNotFound) {
			//Если заказ не найден - проставляем статус INVALID
			err := service.orderRepository.UpdateOrder(ctx, orderID, orderPkg.InvalidStatus, orderEntity.Accrual)
			if err != nil {
				service.logger.Error(err.Error())
				return err
			}
		}
		service.logger.Error(err.Error())
		return err
	}

	status := orderEntity.Status
	switch orderFromAccrualSystem.Status {
	case accrualPkg.AccrualProcessedStatus:
		status = orderPkg.ProcessedStatus
	case accrualPkg.AccrualInvalidStatus:
		status = orderPkg.InvalidStatus
	default:
		status = orderPkg.ProcessingStatus
	}

	err = service.orderRepository.UpdateOrder(ctx, orderID, status, orderFromAccrualSystem.Accrual)
	if err != nil {
		service.logger.Error(err.Error())
		return err
	}

	return nil
}

func (service *OrderService) GetUserOrders(ctx context.Context, userID string) ([]orderPkg.Order, error) {
	orders, err := service.orderRepository.GetOrders(ctx, orderQueryPkg.OrdersQuery{
		Filter: orderQueryPkg.OrdersFilterDto{
			UserID: userID,
		},
		Sort: repository.SortDto{
			By:        orderQueryPkg.ByUploadedAt,
			Direction: repository.AscDirection,
		},
	})
	if err != nil {
		service.logger.Error(err.Error())
		return nil, err
	}

	return orders, nil
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
