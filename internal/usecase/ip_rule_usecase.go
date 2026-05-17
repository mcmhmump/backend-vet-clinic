package usecase

import (
	"net/netip"
	"strings"

	"github.com/mcmhmump/backend-vet-clinic/internal/domain"
)

type IPRuleRepo interface {
	GetAll() ([]domain.IPRule, error)
	Create(rule *domain.IPRule) error
	DeleteByID(id uint) error
}

type IPRuleUsecase struct {
	repo IPRuleRepo
}

func NewIPRuleUsecase(repo IPRuleRepo) *IPRuleUsecase {
	return &IPRuleUsecase{repo: repo}
}

func (u *IPRuleUsecase) GetAll() ([]domain.IPRule, error) {
	return u.repo.GetAll()
}

func (u *IPRuleUsecase) Create(listType, value string) error {
	rule := &domain.IPRule{
		ListType: listType,
		Value:    value,
	}
	return u.repo.Create(rule)
}

func (u *IPRuleUsecase) DeleteByID(id uint) error {
	return u.repo.DeleteByID(id)
}

func (u *IPRuleUsecase) CheckIP(ipStr string) (bool, string, error) {
	rules, err := u.repo.GetAll()
	if err != nil {
		return false, "", err
	}

	clientIP, err := netip.ParseAddr(ipStr)
	if err != nil {
		return false, "invalid_ip_format", nil
	}

	for _, rule := range rules {
		if rule.ListType == "blacklist" && matches(clientIP, rule.Value) {
			return false, "in_blacklist", nil
		}
	}

	for _, rule := range rules {
		if rule.ListType == "whitelist" && matches(clientIP, rule.Value) {
			return true, "in_whitelist", nil
		}
	}

	return false, "default_deny", nil
}

func matches(clientIP netip.Addr, value string) bool {
	if strings.Contains(value, "/") {
		prefix, err := netip.ParsePrefix(value)
		return err == nil && prefix.Contains(clientIP)
	}

	addr, err := netip.ParseAddr(value)
	return err == nil && addr == clientIP
}
