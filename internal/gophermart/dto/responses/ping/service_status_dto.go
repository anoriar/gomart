package ping

import "encoding/json"

type ServiceStatusDto struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Error  string `json:"error"`
}

func (ms ServiceStatusDto) MarshalJSON() ([]byte, error) {
	type alias ServiceStatusDto
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
