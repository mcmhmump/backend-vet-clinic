package http

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/mcmhmump/backend-vet-clinic/internal/usecase"
	"go.uber.org/zap"
)

type ProxyHandler struct {
	reverseProxy *httputil.ReverseProxy
	logger       *zap.Logger
	ipFilter     *usecase.IPFilterService
}

func NewProxyHandler(targetURL string, logger *zap.Logger, ipFilter *usecase.IPFilterService) (*ProxyHandler, error) {
	target, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	return &ProxyHandler{
		reverseProxy: proxy,
		logger:       logger,
		ipFilter:     ipFilter,
	}, nil
}

func (h *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr

	allowed, reason := h.ipFilter.CheckIP(clientIP)
	if !allowed {
		h.logger.Warn("access denied",
			zap.String("client_ip", clientIP),
			zap.String("url", r.URL.Path),
			zap.String("reason", reason),
		)
		http.Error(w, "Forbidden: "+reason, http.StatusForbidden)
		return
	}

	h.logger.Info("access allowed",
		zap.String("client_ip", clientIP),
		zap.String("url", r.URL.Path),
		zap.String("reason", reason),
	)

	h.reverseProxy.ServeHTTP(w, r)
}
