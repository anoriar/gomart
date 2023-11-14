package accrual

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/anoriar/gophermart/internal/gophermart/domain_errors"
	"github.com/anoriar/gophermart/internal/gophermart/dto/accrual"
	"io"
	"net/http"
)

type AccrualRepository struct {
	httpClient *http.Client
	baseURL    string
}

func NewAccrualRepository(httpClient *http.Client, baseURL string) *AccrualRepository {
	return &AccrualRepository{httpClient: httpClient, baseURL: baseURL}
}

func (repository *AccrualRepository) GetOrder(orderId string) (result accrual.AccrualOrderDto, exists bool, err error) {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/orders/%s", repository.baseURL, orderId), nil)
	if err != nil {
		return accrual.AccrualOrderDto{}, false, fmt.Errorf("accrual service. getOrder: %w: %v", domain_errors.ErrInternalError, err)
	}

	response, err := repository.httpClient.Do(request)
	if err != nil {
		return accrual.AccrualOrderDto{}, false, fmt.Errorf("accrual service. getOrder: %w: %v", domain_errors.ErrInternalError, err)
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusOK:
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return accrual.AccrualOrderDto{}, false, fmt.Errorf("accrual service. getOrder: %w: %v", domain_errors.ErrInternalError, err)
		}

		var accrualOrder accrual.AccrualOrderDto

		err = json.Unmarshal(body, &accrualOrder)
		if err != nil {
			return accrual.AccrualOrderDto{}, false, fmt.Errorf("accrual service. getOrder: %w: %v", domain_errors.ErrInternalError, err)
		}

		return accrualOrder, true, nil
	case http.StatusNoContent:
		return accrual.AccrualOrderDto{}, false, nil
	case http.StatusTooManyRequests:
		return accrual.AccrualOrderDto{}, false, fmt.Errorf("accrual service. getOrder: %w", domain_errors.ErrTooManyRequests)
	case http.StatusInternalServerError:
		return accrual.AccrualOrderDto{}, false, fmt.Errorf("accrual service. getOrder: %w: %v", domain_errors.ErrDependencyFailure, errors.New("internal server error"))
	default:
		return accrual.AccrualOrderDto{}, false, fmt.Errorf("accrual service. getOrder: %w: %v", domain_errors.ErrDependencyFailure, errors.New("not expected status"))
	}
}
