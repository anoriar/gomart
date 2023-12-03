package responses

type BalanceResponseDto struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}
