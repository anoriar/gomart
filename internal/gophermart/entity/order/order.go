package order

import "time"

const (
	RegisteredStatus = "REGISTERED"
	ProcessingStatus = "PROCESSING"
	InvalidStatus    = "INVALID"
	ProcessedStatus  = "PROCESSED"
)

type Order struct {
	Id         string    `db:"id"`
	Status     string    `db:"status"`
	Accrual    float64   `db:"accrual"`
	UploadedAt time.Time `db:"uploaded_at"`
	UserID     string    `db:"user_id"`
}

func CreateNewOrder(
	id string,
	userID string,
) Order {
	return Order{
		Id:         id,
		Status:     RegisteredStatus,
		Accrual:    0,
		UploadedAt: time.Now(),
		UserID:     userID,
	}
}
