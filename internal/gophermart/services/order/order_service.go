package order

import (
	"context"
	"errors"
	"fmt"
	"github.com/anoriar/gophermart/internal/gophermart/domain_errors"
	"github.com/anoriar/gophermart/internal/gophermart/dto/repository"
	orderQueryPkg "github.com/anoriar/gophermart/internal/gophermart/dto/repository/order"
	orderPkg "github.com/anoriar/gophermart/internal/gophermart/entity/order"
	"github.com/anoriar/gophermart/internal/gophermart/processors/bus"
	"github.com/anoriar/gophermart/internal/gophermart/processors/order/message"
	"github.com/anoriar/gophermart/internal/gophermart/repository/order"
	balanceServicePkg "github.com/anoriar/gophermart/internal/gophermart/services/balance"
	"github.com/anoriar/gophermart/internal/gophermart/services/order/fetcher"
	"github.com/anoriar/gophermart/internal/gophermart/services/validator/id_validator"
	"go.uber.org/zap"
)

type OrderService struct {
	orderRepository   order.OrderRepositoryInterface
	orderFetchService fetcher.OrderFetchServiceInterface
	balanceService    balanceServicePkg.BalanceServiceInterface
	messageBus        bus.MessageBusInterface
	idValidator       id_validator.IdValidatorInterface
	logger            *zap.Logger
}

func NewOrderService(
	orderRepository order.OrderRepositoryInterface,
	orderFetchService fetcher.OrderFetchServiceInterface,
	balanceService balanceServicePkg.BalanceServiceInterface,
	messageBus bus.MessageBusInterface,
	idValidator id_validator.IdValidatorInterface,
	logger *zap.Logger,
) *OrderService {
	return &OrderService{
		orderRepository:   orderRepository,
		orderFetchService: orderFetchService,
		balanceService:    balanceService,
		idValidator:       idValidator,
		messageBus:        messageBus,
		logger:            logger,
	}
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

	newOrder, err := service.orderFetchService.Fetch(orderEntity)
	if err != nil {
		service.logger.Error(err.Error())
		return err
	}

	err = service.orderRepository.UpdateOrder(ctx, orderID, newOrder.Status, newOrder.Accrual)
	if err != nil {
		service.logger.Error(err.Error())
		return err
	}
	service.logger.Info(fmt.Sprintf("order %s processed successfully", orderID),
		zap.String("status", newOrder.Status),
		zap.Float64("accrual", newOrder.Accrual),
	)

	if newOrder.Accrual > 0 && newOrder.Status == orderPkg.ProcessedStatus {
		err := service.balanceService.UpdateUserBalance(ctx, newOrder.UserID, newOrder.Accrual, 0)
		if err != nil {
			service.logger.Error(err.Error())
			return err
		}
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
	if service.idValidator.Validate(orderID) == false {
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
		err = service.messageBus.SendMessage(message.OrderProcessMessage{OrderID: orderID})
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
