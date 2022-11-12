package repository

import (
	"BIOCAD/internal/structures"
	"gorm.io/gorm"
)

type PaginationPostgres struct {
	db *gorm.DB
}

func NewPaginationPostgres(db *gorm.DB) *PaginationPostgres {
	return &PaginationPostgres{db: db}
}

func (pp *PaginationPostgres) GetAll(guid string, pagination *structures.Pagination) ([]structures.Device, error) {
	var devices []structures.Device
	offset := (pagination.Page - 1) * pagination.Limit
	queryBuider := pp.db.Limit(pagination.Limit).Offset(offset)
	result := queryBuider.Model(&structures.Device{}).Where(`unit_guid = ?`, guid).Find(&devices)
	if result.Error != nil {
		msg := result.Error
		return nil, msg
	}
	return devices, nil
}
