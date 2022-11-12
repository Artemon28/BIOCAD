package repository

import (
	"BIOCAD/internal/structures"
	"gorm.io/gorm"
)

type DevicePostgres struct {
	db *gorm.DB
}

func NewDevicePostgres(db *gorm.DB) *DevicePostgres {
	return &DevicePostgres{db: db}
}

func (r *DevicePostgres) AddDevice(device structures.Device) (int, error) {
	err := r.db.Create(&device).Error

	if err != nil {
		return 0, err
	}
	return device.Id, nil
}

func (r *DevicePostgres) GetDevices(UnitGuid string) ([]structures.Device, error) {
	var devices []structures.Device
	err := r.db.Where(`unit_guid = ?`, UnitGuid).
		Find(&devices).Error
	if err != nil {
		return nil, err
	}
	return devices, nil
}
