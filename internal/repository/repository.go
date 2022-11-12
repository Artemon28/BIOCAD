package repository

import (
	"BIOCAD/internal/structures"
	"gorm.io/gorm"
)

type Device interface {
	AddDevice(device structures.Device) (int, error)
	GetDevices(UnitGuid string) ([]structures.Device, error)
}

type DeviceFile interface {
	AddFile(name string) (int, error)
	GetAllFiles() ([]structures.File, error)
}

type Pagination interface {
	GetAll(guid string, pagination *structures.Pagination) ([]structures.Device, error)
}

type Repository struct {
	Device
	DeviceFile
	Pagination
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Device:     NewDevicePostgres(db),
		DeviceFile: NewDeviceFilePostgres(db),
		Pagination: NewPaginationPostgres(db),
	}
}
