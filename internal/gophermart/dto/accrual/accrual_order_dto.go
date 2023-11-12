package accrual

const (
	AccrualRegisteredStatus = "REGISTERED"
	AccrualProcessingStatus = "PROCESSING"
	AccrualInvalidStatus    = "INVALID"
	AccrualProcessedStatus  = "PROCESSED"
)

type AccrualOrderDto struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}
