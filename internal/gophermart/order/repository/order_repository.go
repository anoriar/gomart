package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	orderRepositoryDtoPkg "github.com/anoriar/gophermart/internal/gophermart/order/dto/repository"
	"github.com/anoriar/gophermart/internal/gophermart/order/entity"
	"github.com/anoriar/gophermart/internal/gophermart/shared/app/db"
	repositoryDtoPkg "github.com/anoriar/gophermart/internal/gophermart/shared/dto/repository"
	errors2 "github.com/anoriar/gophermart/internal/gophermart/shared/errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"strings"
)

type OrderRepository struct {
	db *db.Database
}

func NewOrderRepository(db *db.Database) *OrderRepository {
	return &OrderRepository{db: db}
}

func (repository *OrderRepository) AddOrder(ctx context.Context, order entity.Order) error {
	_, err := repository.db.Conn.ExecContext(ctx,
		`INSERT INTO orders (id, status, accrual, uploaded_at, user_id) 
							VALUES ($1, $2, $3, $4, $5)`,
		order.ID, order.Status, order.Accrual, order.UploadedAt, order.UserID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
			return fmt.Errorf("%w: %v", errors2.ErrConflict, err)
		}
		return fmt.Errorf("%w: %v", errors2.ErrInternalError, err)
	}
	return nil
}

func (repository *OrderRepository) GetOrderByID(ctx context.Context, orderID string) (entity.Order, error) {
	var orderRes entity.Order
	err := repository.db.Conn.QueryRowxContext(ctx, "SELECT id, status, accrual, uploaded_at, user_id FROM orders WHERE id=$1", orderID).StructScan(&orderRes)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return orderRes, fmt.Errorf("%w: %v", errors2.ErrNotFound, err)
		}
		return orderRes, fmt.Errorf("%w: %v", errors2.ErrInternalError, err)
	}
	return orderRes, nil
}

func (repository *OrderRepository) UpdateOrder(ctx context.Context, orderID string, status string, accrual float64) error {
	_, err := repository.db.Conn.NamedExecContext(ctx, "UPDATE orders SET status = :status, accrual = :accrual WHERE id = :id", map[string]interface{}{
		"id":      orderID,
		"status":  status,
		"accrual": accrual,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: %v", errors2.ErrNotFound, err)
		}
		return fmt.Errorf("%w: %v", errors2.ErrInternalError, err)
	}

	return nil
}

func (repository *OrderRepository) buildSort(queryRowSlice *[]string, sortDto repositoryDtoPkg.SortDto) {
	switch sortDto.By {
	case orderRepositoryDtoPkg.ByUploadedAt:
		*queryRowSlice = append(*queryRowSlice, "ORDER BY uploaded_at")
	default:
		*queryRowSlice = append(*queryRowSlice, "ORDER BY id")
	}
	switch sortDto.Direction {
	case repositoryDtoPkg.DescDirection:
		*queryRowSlice = append(*queryRowSlice, "DESC")
	default:
		*queryRowSlice = append(*queryRowSlice, "ASC")
	}
}

func (repository *OrderRepository) buildFilter(queryRowSlice *[]string, queryParams *[]interface{}, filterDto orderRepositoryDtoPkg.OrdersFilterDto) {
	var filters []string
	if filterDto.UserID != "" {
		filters = append(filters, "user_id = ?")
		*queryParams = append(*queryParams, filterDto.UserID)
	}
	if len(filterDto.Statuses) > 0 {

		var statusPlaceholders []string
		for _, statusStr := range filterDto.Statuses {
			statusPlaceholders = append(statusPlaceholders, "?")
			*queryParams = append(*queryParams, statusStr)
		}
		filters = append(filters, fmt.Sprintf("status IN (%s)", strings.Join(statusPlaceholders, ", ")))
	}

	if len(filters) != 0 {
		*queryRowSlice = append(*queryRowSlice, "WHERE "+strings.Join(filters, " AND "))
	}
}

func (repository *OrderRepository) buildPagination(queryRowSlice *[]string, queryParams *[]interface{}, paginationDto repositoryDtoPkg.PaginationDto) {
	*queryRowSlice = append(*queryRowSlice, "OFFSET ?")
	*queryParams = append(*queryParams, paginationDto.Offset)

	if paginationDto.Limit != 0 {
		*queryRowSlice = append(*queryRowSlice, "LIMIT ?")
		*queryParams = append(*queryParams, paginationDto.Limit)
	}
}

func (repository *OrderRepository) GetOrders(ctx context.Context, query orderRepositoryDtoPkg.OrdersQuery) ([]entity.Order, error) {
	var orders []entity.Order
	var queryParams []interface{}
	var queryRowSlice []string
	queryRowSlice = append(queryRowSlice, "SELECT id, user_id, status, accrual, uploaded_at FROM orders")

	repository.buildFilter(&queryRowSlice, &queryParams, query.Filter)
	repository.buildSort(&queryRowSlice, query.Sort)
	repository.buildPagination(&queryRowSlice, &queryParams, query.Pagination)

	queryRow := strings.Join(queryRowSlice, " ")
	queryRow = repository.db.Conn.Rebind(queryRow)
	rows, err := repository.db.Conn.QueryContext(ctx, queryRow, queryParams...)
	if err != nil {
		return nil, fmt.Errorf("order repository GetOrders: %w: %v", errors2.ErrInternalError, err)
	}
	defer rows.Close()

	for rows.Next() {
		var order entity.Order
		err := rows.Scan(&order.ID, &order.UserID, &order.Status, &order.Accrual, &order.UploadedAt)
		if err != nil {
			return nil, fmt.Errorf("order repository GetOrders: %w: %v", errors2.ErrInternalError, err)
		}

		orders = append(orders, order)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("order repository GetOrders: %w: %v", errors2.ErrInternalError, rows.Err())
	}

	return orders, nil
}

func (repository *OrderRepository) GetTotal(ctx context.Context, filter orderRepositoryDtoPkg.OrdersFilterDto) (int, error) {
	var count int
	var queryParams []interface{}
	var queryRowSlice []string
	queryRowSlice = append(queryRowSlice, "SELECT count(*) FROM orders")

	repository.buildFilter(&queryRowSlice, &queryParams, filter)

	queryRow := strings.Join(queryRowSlice, " ")
	queryRow = repository.db.Conn.Rebind(queryRow)

	err := repository.db.Conn.GetContext(ctx, &count, queryRow, queryParams...)
	if err != nil {
		return 0, err
	}

	return count, nil
}
