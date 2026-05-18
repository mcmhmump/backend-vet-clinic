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

// GetMetrics godoc
// @Summary Получить метрики
// @Description Возвращает статистику сервиса в памяти
// @Tags monitoring
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /example/metrics [get]
func (h *MonitoringHandler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.metrics.GetMetrics())
}
