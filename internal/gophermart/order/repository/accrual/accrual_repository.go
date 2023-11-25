package accrual

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/anoriar/gophermart/internal/gophermart/order/dto/accrual"
	errors2 "github.com/anoriar/gophermart/internal/gophermart/shared/errors"
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

func (repository *AccrualRepository) GetOrder(orderID string) (result accrual.AccrualOrderDto, exists bool, err error) {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/orders/%s", repository.baseURL, orderID), nil)
	if err != nil {
		return accrual.AccrualOrderDto{}, false, fmt.Errorf("accrual service. getOrder: %w: %v", errors2.ErrInternalError, err)
	}

	response, err := repository.httpClient.Do(request)
	if err != nil {
		return accrual.AccrualOrderDto{}, false, fmt.Errorf("accrual service. getOrder: %w: %v", errors2.ErrInternalError, err)
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusOK:
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return accrual.AccrualOrderDto{}, false, fmt.Errorf("accrual service. getOrder: %w: %v", errors2.ErrInternalError, err)
		}

		var accrualOrder accrual.AccrualOrderDto

		err = json.Unmarshal(body, &accrualOrder)
		if err != nil {
			return accrual.AccrualOrderDto{}, false, fmt.Errorf("accrual service. getOrder: %w: %v", errors2.ErrInternalError, err)
		}

		return accrualOrder, true, nil
	case http.StatusNoContent:
		return accrual.AccrualOrderDto{}, false, nil
	case http.StatusTooManyRequests:
		return accrual.AccrualOrderDto{}, false, fmt.Errorf("accrual service. getOrder: %w", errors2.ErrTooManyRequests)
	case http.StatusInternalServerError:
		return accrual.AccrualOrderDto{}, false, fmt.Errorf("accrual service. getOrder: %w: %v", errors2.ErrDependencyFailure, errors.New("internal server error"))
	default:
		return accrual.AccrualOrderDto{}, false, fmt.Errorf("accrual service. getOrder: %w: %v", errors2.ErrDependencyFailure, errors.New("not expected status"))
	}
}
