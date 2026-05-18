package repository

import (
	"github.com/mcmhmump/backend-vet-clinic/internal/domain"
	"gorm.io/gorm"
)

type AccessLogRepository struct {
	db *gorm.DB
}

func NewAccessLogRepository(db *gorm.DB) *AccessLogRepository {
	return &AccessLogRepository{db: db}
}

func (r *AccessLogRepository) Create(log *domain.AccessLog) error {
	return r.db.Create(log).Error
}
