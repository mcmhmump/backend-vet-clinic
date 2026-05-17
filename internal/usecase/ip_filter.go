package usecase

import (
	"net"
	"net/netip"
	"strings"

	"github.com/mcmhmump/backend-vet-clinic/internal/domain"
)

type IPFilterService struct {
	config *domain.IPAccessConfig
}

func NewIPFilterService(cfg *domain.IPAccessConfig) *IPFilterService {
	return &IPFilterService{config: cfg}
}

func (s *IPFilterService) CheckIP(ipStr string) (bool, string) {
	host := ipStr

	if strings.Contains(ipStr, ":") {
		if h, _, err := net.SplitHostPort(ipStr); err == nil {
			host = h
		}
	}

	clientIP, err := netip.ParseAddr(host)
	if err != nil {
		return false, "invalid_ip_format"
	}

	if s.isInList(clientIP, s.config.Blacklist) {
		return false, "in_blacklist"
	}

	if s.isInList(clientIP, s.config.Whitelist) {
		return true, "in_whitelist"
	}

	if s.config.DefaultPolicy == "allow" {
		return true, "default_allow"
	}

	return false, "default_deny"
}

func (s *IPFilterService) isInList(clientIP netip.Addr, list []string) bool {
	for _, item := range list {
		if strings.Contains(item, "/") {
			prefix, err := netip.ParsePrefix(item)
			if err == nil && prefix.Contains(clientIP) {
				return true
			}
			continue
		}

		addr, err := netip.ParseAddr(item)
		if err == nil && addr.Compare(clientIP) == 0 {
			return true
		}
	}

	return false
}
