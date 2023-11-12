package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/anoriar/gophermart/internal/gophermart/app/db"
	"github.com/anoriar/gophermart/internal/gophermart/domain_errors"
	repository2 "github.com/anoriar/gophermart/internal/gophermart/dto/repository"
	orderQueryPkg "github.com/anoriar/gophermart/internal/gophermart/dto/repository/order"
	"github.com/anoriar/gophermart/internal/gophermart/entity/order"
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

func (repository *OrderRepository) AddOrder(ctx context.Context, order order.Order) error {
	_, err := repository.db.Conn.ExecContext(ctx,
		`INSERT INTO orders (id, status, accrual, uploaded_at, user_id) 
							VALUES ($1, $2, $3, $4, $5)`,
		order.Id, order.Status, order.Accrual, order.UploadedAt, order.UserID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
			return fmt.Errorf("%w: %v", domain_errors.ErrConflict, err)
		}
		return fmt.Errorf("%w: %v", domain_errors.ErrInternalError, err)
	}
	return nil
}

func (repository *OrderRepository) GetOrderByID(ctx context.Context, orderID string) (order.Order, error) {
	var orderRes order.Order
	err := repository.db.Conn.QueryRowxContext(ctx, "SELECT id, status, accrual, uploaded_at, user_id FROM orders WHERE id=$1", orderID).StructScan(&orderRes)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return orderRes, fmt.Errorf("%w: %v", domain_errors.ErrNotFound, err)
		}
		return orderRes, fmt.Errorf("%w: %v", domain_errors.ErrInternalError, err)
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
			return fmt.Errorf("%w: %v", domain_errors.ErrNotFound, err)
		}
		return fmt.Errorf("%w: %v", domain_errors.ErrInternalError, err)
	}

	return nil
}

func (repository *OrderRepository) buildSort(queryRowSlice []string, sortDto repository2.SortDto) {
	switch sortDto.By {
	case orderQueryPkg.ByUploadedAt:
		queryRowSlice = append(queryRowSlice, "ORDER BY uploaded_at")
	default:
		queryRowSlice = append(queryRowSlice, "ORDER BY id")
	}
	switch sortDto.Direction {
	case repository2.DescDirection:
		queryRowSlice = append(queryRowSlice, "DESC")
	default:
		queryRowSlice = append(queryRowSlice, "ASC")
	}
}

func (repository *OrderRepository) buildFilter(queryRowSlice []string, queryParams map[string]interface{}, filterDto orderQueryPkg.OrdersFilterDto) {
	var filters []string
	if filterDto.UserID != "" {
		filters = append(filters, "user_id = :userID")
		queryParams["userID"] = filterDto.UserID
	}
	if filterDto.Status != "" {
		filters = append(filters, "status = :status")
		queryParams["status"] = filterDto.Status
	}

	if len(filters) != 0 {
		queryRowSlice = append(queryRowSlice, "WHERE "+strings.Join(filters, " AND "))
	}
}

func (repository *OrderRepository) buildPagination(queryRowSlice []string, queryParams map[string]interface{}, paginationDto repository2.PaginationDto) {
	queryRowSlice = append(queryRowSlice, "OFFSET :offset")
	queryParams["offset"] = paginationDto.Offset

	if paginationDto.Limit != 0 {
		queryRowSlice = append(queryRowSlice, "LIMIT :limit")
		queryParams["limit"] = paginationDto.Limit
	}
}

func (repository *OrderRepository) GetOrders(ctx context.Context, query orderQueryPkg.OrdersQuery) ([]order.Order, error) {
	var orders []order.Order
	queryParams := make(map[string]interface{})
	var queryRowSlice []string
	queryRowSlice = append(queryRowSlice, "SELECT id, user_id, status, accrual, uploaded_at FROM orders")

	repository.buildFilter(queryRowSlice, queryParams, query.Filter)
	repository.buildSort(queryRowSlice, query.Sort)
	repository.buildPagination(queryRowSlice, queryParams, query.Pagination)

	queryRow := strings.Join(queryRowSlice, " ")
	rows, err := repository.db.Conn.NamedQueryContext(ctx, queryRow, queryParams)
	if err != nil {
		return nil, fmt.Errorf("order repository GetOrders: %w: %v", domain_errors.ErrInternalError, err)
	}
	defer rows.Close()

	for rows.Next() {
		var order order.Order
		err := rows.Scan(&order.Id, &order.UserID, &order.Status, &order.Accrual, &order.UploadedAt)
		if err != nil {
			return nil, fmt.Errorf("order repository GetOrders: %w: %v", domain_errors.ErrInternalError, err)
		}

		orders = append(orders, order)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("order repository GetOrders: %w: %v", domain_errors.ErrInternalError, rows.Err())
	}

	return orders, nil
}

func (repository *OrderRepository) GetTotal(ctx context.Context, filter orderQueryPkg.OrdersFilterDto) (int, error) {

	//TODO implement me
	panic("implement me")
}
