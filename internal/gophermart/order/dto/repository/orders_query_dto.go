package repository

import (
	repository2 "github.com/anoriar/gophermart/internal/gophermart/shared/dto/repository"
)

const (
	ByUploadedAt = "uploaded_at"
)

type OrdersQuery struct {
	Filter     OrdersFilterDto
	Pagination repository2.PaginationDto
	Sort       repository2.SortDto
}
