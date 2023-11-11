package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/anoriar/gophermart/internal/gophermart/app/db"
	"github.com/anoriar/gophermart/internal/gophermart/domain_errors"
	"github.com/anoriar/gophermart/internal/gophermart/entity/order"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
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
