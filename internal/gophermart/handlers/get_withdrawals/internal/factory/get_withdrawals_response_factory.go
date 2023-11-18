package factory

import (
	withdrawalsResponseDtoPkg "github.com/anoriar/gophermart/internal/gophermart/dto/responses/withdrawals"
	"github.com/anoriar/gophermart/internal/gophermart/entity/withdrawal"
	"time"
)

type GetWithdrawalsResponseFactory struct {
}

func NewGetWithdrawalsResponseFactory() *GetWithdrawalsResponseFactory {
	return &GetWithdrawalsResponseFactory{}
}

func (factory GetWithdrawalsResponseFactory) CreateWithdrawalsResponse(withdrawals []withdrawal.Withdrawal) []withdrawalsResponseDtoPkg.WithdrawalResponseDto {
	var response []withdrawalsResponseDtoPkg.WithdrawalResponseDto
	for _, withdrawalEntity := range withdrawals {
		response = append(response, withdrawalsResponseDtoPkg.WithdrawalResponseDto{
			Order:       withdrawalEntity.OrderID,
			Sum:         withdrawalEntity.Sum,
			ProcessedAt: withdrawalEntity.ProcessedAt.Format(time.RFC3339),
		})
	}
	return response
}
