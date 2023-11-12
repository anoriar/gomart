package order

import "github.com/anoriar/gophermart/internal/gophermart/dto/repository"

const (
	ByUploadedAt = "uploaded_at"
)

type OrdersQuery struct {
	Filter     OrdersFilterDto
	Pagination repository.PaginationDto
	Sort       repository.SortDto
}
