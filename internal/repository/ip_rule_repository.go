package repository

import (
	"github.com/mcmhmump/backend-vet-clinic/internal/domain"
	"gorm.io/gorm"
)

type IPRuleRepository struct {
	db *gorm.DB
}

func NewIPRuleRepository(db *gorm.DB) *IPRuleRepository {
	return &IPRuleRepository{db: db}
}

func (r *IPRuleRepository) GetAll() ([]domain.IPRule, error) {
	var rules []domain.IPRule
	err := r.db.Find(&rules).Error
	return rules, err
}

func (r *IPRuleRepository) Create(rule *domain.IPRule) error {
	return r.db.Create(rule).Error
}

func (r *IPRuleRepository) DeleteByID(id uint) error {
	return r.db.Delete(&domain.IPRule{}, id).Error
}
