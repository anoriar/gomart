package ping

const (
	OKStatus   = "OK"
	FailStatus = "FAIL"
)

type PingResponseDto struct {
	Services []ServiceStatusDto `json:"services"`
}
