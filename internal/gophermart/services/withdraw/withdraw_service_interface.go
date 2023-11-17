package withdraw

import (
	"context"
	"github.com/anoriar/gophermart/internal/gophermart/dto/requests/withdraw"
)

type WithdrawServiceInterface interface {
	Withdraw(ctx context.Context, userID string, withdraw withdraw.WithdrawDto) error
}
