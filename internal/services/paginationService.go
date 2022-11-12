package services

import (
	"BIOCAD/internal/repository"
	"BIOCAD/internal/structures"
)

type PaginationService struct {
	r *repository.Repository
}

func NewPaginationService(r *repository.Repository) *PaginationService {
	return &PaginationService{r: r}
}

func (ps *PaginationService) GetAll(guid string, pagination *structures.Pagination) ([]structures.Device, error) {
	return ps.r.GetAll(guid, pagination)
}
