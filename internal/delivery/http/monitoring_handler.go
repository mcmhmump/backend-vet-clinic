package http

import (
	"encoding/json"
	"net/http"

	"github.com/mcmhmump/backend-vet-clinic/internal/usecase"
)

type MonitoringHandler struct {
	metrics *usecase.MetricsService
}

func NewMonitoringHandler(metrics *usecase.MetricsService) *MonitoringHandler {
	return &MonitoringHandler{metrics: metrics}
}

func (h *MonitoringHandler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.metrics.GetMetrics())
}
