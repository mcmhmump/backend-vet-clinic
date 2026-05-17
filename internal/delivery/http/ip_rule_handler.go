package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mcmhmump/backend-vet-clinic/internal/usecase"
	"go.uber.org/zap"
)

type IPRuleHandler struct {
	usecase *usecase.IPRuleUsecase
	logger  *zap.Logger
}

type CreateIPRuleRequest struct {
	ListType string `json:"list_type"`
	Value    string `json:"value"`
}

func NewIPRuleHandler(usecase *usecase.IPRuleUsecase, logger *zap.Logger) *IPRuleHandler {
	return &IPRuleHandler{
		usecase: usecase,
		logger:  logger,
	}
}

func (h *IPRuleHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	rules, err := h.usecase.GetAll()
	if err != nil {
		http.Error(w, "failed to load rules", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rules)
}

func (h *IPRuleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateIPRuleRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.ListType == "" || req.Value == "" {
		http.Error(w, "list_type and value are required", http.StatusBadRequest)
		return
	}

	if err := h.usecase.Create(req.ListType, req.Value); err != nil {
		http.Error(w, "failed to create rule", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "created",
	})
}

func (h *IPRuleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.usecase.DeleteByID(uint(id64)); err != nil {
		http.Error(w, "failed to delete rule", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "deleted",
	})
}

func (h *IPRuleHandler) CheckIP(w http.ResponseWriter, r *http.Request) {
	ip := r.URL.Query().Get("ip")
	if ip == "" {
		http.Error(w, "ip query param is required", http.StatusBadRequest)
		return
	}

	allowed, reason, err := h.usecase.CheckIP(ip)
	if err != nil {
		http.Error(w, "failed to check ip", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ip":      ip,
		"allowed": allowed,
		"reason":  reason,
	})
}
