package services

import (
	"BIOCAD/internal/repository"
	"BIOCAD/internal/structures"
	"sync"
	"time"
)

type DirScanInterface interface {
	Scan(dirName string, duration time.Duration)
	MakeReports(unitGuidChan chan structures.Device)
	MakeReportFile(unitGuid structures.Device, wg *sync.WaitGroup)
}

type PaginationInterface interface {
	GetAll(guid string, pagination *structures.Pagination) ([]structures.Device, error)
}

type Service struct {
	DirScanInterface
	PaginationInterface
}

func NewService(rep *repository.Repository) *Service {
	return &Service{
		DirScanInterface:    NewScanDirectory(rep),
		PaginationInterface: NewPaginationService(rep),
	}
}
