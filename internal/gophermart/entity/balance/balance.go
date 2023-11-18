package balance

import "time"

type Balance struct {
	ID         string    `db:"id"`
	UserID     string    `db:"user_id"`
	Balance    float64   `db:"balance"`
	Withdrawal float64   `db:"withdrawal"`
	UpdatedAt  time.Time `db:"updated_at"`
}
