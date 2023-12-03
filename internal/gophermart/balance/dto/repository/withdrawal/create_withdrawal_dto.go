package withdrawal

type CreateWithdrawalDto struct {
	UserID string
	Order  string
	Sum    float64
}
