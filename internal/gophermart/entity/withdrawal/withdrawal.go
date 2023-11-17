package withdrawal

import "time"

type Withdrawal struct {
	Id          string    `db:"id"`
	UserID      string    `db:"user_id"`
	OrderID     string    `db:"order_id"`
	Sum         float64   `db:"sum"`
	ProcessedAt time.Time `db:"processed_at"`
}
