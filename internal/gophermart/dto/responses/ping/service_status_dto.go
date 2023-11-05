package ping

import "encoding/json"

type ServiceStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Error  string `json:"error"`
}

func (ms ServiceStatus) MarshalJSON() ([]byte, error) {
	type alias ServiceStatus
	if ms.Error == "" {
		return json.Marshal(&struct {
			Name   string      `json:"name"`
			Status string      `json:"status"`
			Error  interface{} `json:"error"`
		}{
			Name:   ms.Name,
			Status: ms.Status,
			Error:  nil,
		})
	}

	return json.Marshal((alias)(ms))
}
