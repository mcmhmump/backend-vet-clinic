package repository

import (
	"github.com/glebarez/sqlite"
	"github.com/mcmhmump/backend-vet-clinic/internal/domain"
	"gorm.io/gorm"
)

func NewDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("vetclinic_proxy.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&domain.IPRule{},
		&domain.AccessLog{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
