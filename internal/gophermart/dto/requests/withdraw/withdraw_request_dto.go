package withdraw

type WithdrawDto struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}
