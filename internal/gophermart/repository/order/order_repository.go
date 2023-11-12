package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/anoriar/gophermart/internal/gophermart/app/db"
	"github.com/anoriar/gophermart/internal/gophermart/domain_errors"
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

func (repository *OrderRepository) GetOrders(ctx context.Context, query orderQueryPkg.OrdersQuery) ([]order.Order, error) {
	var orders []order.Order
	var filters []string
	queryParams := make(map[string]interface{})
	queryRow := "SELECT id, user_id, status, accrual, uploaded_at FROM orders"

	if query.Filter.UserID != "" {
		filters = append(filters, "user_id = :userID")
		queryParams["userID"] = query.Filter.UserID
	}
	if query.Filter.Status != "" {
		filters = append(filters, "status = :status")
		queryParams["status"] = query.Filter.Status
	}

	if len(filters) != 0 {
		queryRow += " WHERE " + strings.Join(filters, " AND ")
	}

	if query.Pagination.Limit != 0 {
		queryRow += " LIMIT :limit"
		queryParams["limit"] = query.Pagination.Limit
	}

	queryRow += " OFFSET :offset"
	queryParams["offset"] = query.Pagination.Offset
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
