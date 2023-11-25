package factory

import (
	"github.com/anoriar/gophermart/internal/gophermart/balance/dto/requests"
	"github.com/anoriar/gophermart/internal/gophermart/balance/entity"
	"github.com/google/uuid"
	"time"
)

type WithdrawalFactory struct {
}

func NewWithdrawalFactory() *WithdrawalFactory {
	return &WithdrawalFactory{}
}

func (WithdrawalFactory) CreateEntityFromRequest(userID string, withdrawDto requests.WithdrawDto) entity.Withdrawal {
	return entity.Withdrawal{
		ID:          uuid.NewString(),
		UserID:      userID,
		OrderID:     withdrawDto.Order,
		Sum:         withdrawDto.Sum,
		ProcessedAt: time.Now(),
	}
}
