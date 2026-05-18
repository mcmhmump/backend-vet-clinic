package usecase

import (
	"testing"

	"github.com/mcmhmump/backend-vet-clinic/internal/domain"
	"github.com/stretchr/testify/assert"
)

type mockIPRuleRepo struct {
	rules []domain.IPRule
}

func (m *mockIPRuleRepo) GetAll() ([]domain.IPRule, error) {
	return m.rules, nil
}

func (m *mockIPRuleRepo) Create(rule *domain.IPRule) error {
	return nil
}

func (m *mockIPRuleRepo) DeleteByID(id uint) error {
	return nil
}

func TestIPRuleUsecase_CheckIP(t *testing.T) {
	tests := []struct {
		name    string
		ip      string
		rules   []domain.IPRule
		allowed bool
		reason  string
	}{
		{
			name: "whitelist allows ip",
			ip:   "127.0.0.1",
			rules: []domain.IPRule{
				{ListType: "whitelist", Value: "127.0.0.1"},
			},
			allowed: true,
			reason:  "in_whitelist",
		},
		{
			name: "blacklist blocks ip",
			ip:   "127.0.0.1",
			rules: []domain.IPRule{
				{ListType: "blacklist", Value: "127.0.0.1"},
			},
			allowed: false,
			reason:  "in_blacklist",
		},
		{
			name: "blacklist has priority",
			ip:   "127.0.0.1",
			rules: []domain.IPRule{
				{ListType: "whitelist", Value: "127.0.0.1"},
				{ListType: "blacklist", Value: "127.0.0.1"},
			},
			allowed: false,
			reason:  "in_blacklist",
		},
		{
			name: "default deny",
			ip:   "10.10.10.10",
			rules: []domain.IPRule{
				{ListType: "whitelist", Value: "127.0.0.1"},
			},
			allowed: false,
			reason:  "default_deny",
		},
		{
			name:    "invalid ip format",
			ip:      "bad-ip",
			rules:   []domain.IPRule{},
			allowed: false,
			reason:  "invalid_ip_format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockIPRuleRepo{rules: tt.rules}
			uc := NewIPRuleUsecase(repo)

			allowed, reason, err := uc.CheckIP(tt.ip)

			assert.NoError(t, err)
			assert.Equal(t, tt.allowed, allowed)
			assert.Equal(t, tt.reason, reason)
		})
	}
}
