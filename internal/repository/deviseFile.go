package repository

import (
	"BIOCAD/internal/structures"
	"gorm.io/gorm"
)

type DeviceFilePostgres struct {
	db *gorm.DB
}

func NewDeviceFilePostgres(db *gorm.DB) *DeviceFilePostgres {
	return &DeviceFilePostgres{db: db}
}

func (r *DeviceFilePostgres) AddFile(name string) (int, error) {
	file := structures.File{Name: name}
	err := r.db.Create(&file).Error

	if err != nil {
		return 0, err
	}
	return file.Id, nil
}

func (r *DeviceFilePostgres) GetAllFiles() ([]structures.File, error) {
	var files []structures.File
	err := r.db.Find(&files).Error
	if err != nil {
		return nil, err
	}
	return files, nil
}
