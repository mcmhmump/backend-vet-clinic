package usecase

import (
	"sync"
)

// MetricsService хранит статистику в памяти (Требование 3.2.1.3)
type MetricsService struct {
	mu              sync.RWMutex
	TotalRequests   int64 `json:"total_requests"`
	AllowedRequests int64 `json:"allowed_requests"`
	BlockedRequests int64 `json:"blocked_requests"`
}

func NewMetricsService() *MetricsService {
	return &MetricsService{}
}

func (m *MetricsService) RecordRequest(allowed bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.TotalRequests++
	if allowed {
		m.AllowedRequests++
	} else {
		m.BlockedRequests++
	}
}

func (m *MetricsService) GetMetrics() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return map[string]interface{}{
		"total_requests":   m.TotalRequests,
		"allowed_requests": m.AllowedRequests,
		"blocked_requests": m.BlockedRequests,
	}
}
