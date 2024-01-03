package services

import (
	"context"
	"errors"
	"fmt"
	balanceServicePkg "github.com/anoriar/gophermart/internal/gophermart/balance/services/balance"
	orderRepositoryDtoPkg "github.com/anoriar/gophermart/internal/gophermart/order/dto/repository"
	orderPkg "github.com/anoriar/gophermart/internal/gophermart/order/entity"
	errors2 "github.com/anoriar/gophermart/internal/gophermart/order/errors"
	"github.com/anoriar/gophermart/internal/gophermart/order/processors/orderprocess/message"
	orderRepositoryPkg "github.com/anoriar/gophermart/internal/gophermart/order/repository"
	"github.com/anoriar/gophermart/internal/gophermart/order/services/fetcher"
	"github.com/anoriar/gophermart/internal/gophermart/shared/dto/repository"
	errors3 "github.com/anoriar/gophermart/internal/gophermart/shared/errors"
	"github.com/anoriar/gophermart/internal/gophermart/shared/services/bus"
	"github.com/anoriar/gophermart/internal/gophermart/shared/services/validator/idvalidator"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type OrderService struct {
	orderRepository   orderRepositoryPkg.OrderRepositoryInterface
	orderFetchService fetcher.OrderFetchServiceInterface
	balanceService    balanceServicePkg.BalanceServiceInterface
	messageBus        bus.MessageBusInterface
	idValidator       idvalidator.IDValidatorInterface
	logger            *zap.Logger
}

func NewOrderService(
	orderRepository orderRepositoryPkg.OrderRepositoryInterface,
	orderFetchService fetcher.OrderFetchServiceInterface,
	balanceService balanceServicePkg.BalanceServiceInterface,
	messageBus bus.MessageBusInterface,
	idValidator idvalidator.IDValidatorInterface,
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderService::ProcessOrder")
	span.SetTag("orderId", orderID)
	defer span.Finish()

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

	newOrder, err := service.orderFetchService.Fetch(ctx, orderEntity)
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
		err := service.balanceService.AddUserBalance(ctx, newOrder.UserID, newOrder.Accrual)
		if err != nil {
			service.logger.Error(err.Error())
			return err
		}
	}

	return nil
}

func (service *OrderService) GetUserOrders(ctx context.Context, userID string) ([]orderPkg.Order, error) {
	orders, err := service.orderRepository.GetOrders(ctx, orderRepositoryDtoPkg.OrdersQuery{
		Filter: orderRepositoryDtoPkg.OrdersFilterDto{
			UserID: userID,
		},
		Sort: repository.SortDto{
			By:        orderRepositoryDtoPkg.ByUploadedAt,
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
	if !service.idValidator.Validate(orderID) {
		return errors2.ErrOrderNumberNotValid
	}
	currentOrder, err := service.orderRepository.GetOrderByID(ctx, orderID)

	switch {
	case err != nil && errors.Is(err, errors3.ErrNotFound):
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
	case err != nil && !errors.Is(err, errors3.ErrNotFound):
		service.logger.Error(err.Error())
		return err
	case currentOrder.UserID == userID:
		return errors2.ErrOrderAlreadyLoaded
	case currentOrder.UserID != userID:
		return errors2.ErrOrderLoadedByAnotherUser
	default:
		service.logger.Error("unexpected behaviour")
		return fmt.Errorf("unexpected behaviour %w", errors3.ErrInternalError)
	}
}

func (service *OrderService) GetOrders(ctx context.Context, query orderRepositoryDtoPkg.OrdersQuery) ([]orderPkg.Order, error) {
	return service.orderRepository.GetOrders(ctx, query)
}

func (service *OrderService) GetOrdersTotal(ctx context.Context, filter orderRepositoryDtoPkg.OrdersFilterDto) (int, error) {
	return service.orderRepository.GetTotal(ctx, filter)
}
