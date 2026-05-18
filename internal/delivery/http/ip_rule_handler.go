package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/mcmhmump/backend-vet-clinic/internal/domain"
	"github.com/mcmhmump/backend-vet-clinic/internal/repository"
	"github.com/mcmhmump/backend-vet-clinic/internal/usecase"
	"go.uber.org/zap"
)

type IPRuleHandler struct {
	usecase *usecase.IPRuleUsecase
	logger  *zap.Logger
	logRepo *repository.AccessLogRepository
	metrics *usecase.MetricsService
}

type CreateIPRuleRequest struct {
	ListType string `json:"list_type"`
	Value    string `json:"value"`
}

func NewIPRuleHandler(u *usecase.IPRuleUsecase, l *zap.Logger, repo *repository.AccessLogRepository, m *usecase.MetricsService) *IPRuleHandler {
	return &IPRuleHandler{usecase: u, logger: l, logRepo: repo, metrics: m}
}

// GetAll godoc
// @Summary Получить все IP-правила
// @Description Возвращает список правил доступа из базы данных
// @Tags ip_access
// @Produce json
// @Success 200 {array} domain.IPRule
// @Router /example/ip_access/allowlists [get]
func (h *IPRuleHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	rules, _ := h.usecase.GetAll()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rules)
}

// Create godoc
// @Summary Добавить IP-правило
// @Description Создает новое правило в whitelist/blacklist/greylist
// @Tags ip_access
// @Accept json
// @Produce json
// @Param request body CreateIPRuleRequest true "Новое IP-правило"
// @Success 201 {object} map[string]string
// @Router /example/ip_access/allowlists [post]
func (h *IPRuleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateIPRuleRequest
	json.NewDecoder(r.Body).Decode(&req)
	h.usecase.Create(req.ListType, req.Value)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "created"})
}

// Delete godoc
// @Summary Удалить IP-правило
// @Description Удаляет правило по ID
// @Tags ip_access
// @Produce json
// @Param id path int true "ID правила"
// @Success 200 {object} map[string]string
// @Router /example/ip_access/allowlists/{id} [delete]
func (h *IPRuleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id64, _ := strconv.ParseUint(idStr, 10, 64)
	h.usecase.DeleteByID(uint(id64))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}

// CheckIP godoc
// @Summary Проверить IP
// @Description Проверяет, разрешен ли доступ для указанного IP
// @Tags ip_access
// @Produce json
// @Param ip query string true "IP адрес"
// @Success 200 {object} map[string]interface{}
// @Router /example/ip_access/check [get]
func (h *IPRuleHandler) CheckIP(w http.ResponseWriter, r *http.Request) {
	ip := r.URL.Query().Get("ip")
	if ip == "" {
		http.Error(w, "ip is required", http.StatusBadRequest)
		return
	}

	allowed, reason, _ := h.usecase.CheckIP(ip)

	// АСИНХРОННАЯ ЗАПИСЬ (Требование 3.1.3)
	go func() {
		logEntry := &domain.AccessLog{
			ClientIP:  ip,
			URL:       r.URL.Path,
			Allowed:   allowed,
			Reason:    reason,
			CreatedAt: time.Now(),
		}
		h.logRepo.Create(logEntry)       // Пишем в SQLite
		h.metrics.RecordRequest(allowed) // Обновляем метрики в памяти
	}()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ip":      ip,
		"allowed": allowed,
		"reason":  reason,
	})
}
