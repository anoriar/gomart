package factory

import (
	"github.com/anoriar/gophermart/internal/gophermart/dto/requests/withdraw"
	"github.com/anoriar/gophermart/internal/gophermart/entity/withdrawal"
	"github.com/google/uuid"
	"time"
)

type WithdrawalFactory struct {
}

func NewWithdrawalFactory() *WithdrawalFactory {
	return &WithdrawalFactory{}
}

func (WithdrawalFactory) CreateEntityFromRequest(userID string, withdrawDto withdraw.WithdrawDto) withdrawal.Withdrawal {
	return withdrawal.Withdrawal{
		Id:          uuid.NewString(),
		UserID:      userID,
		OrderID:     withdrawDto.Order,
		Sum:         withdrawDto.Sum,
		ProcessedAt: time.Now(),
	}
}
