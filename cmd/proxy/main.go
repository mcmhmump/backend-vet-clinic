package main

import (
	"encoding/json"
	"net/http"
	"time"

	deliveryHttp "github.com/mcmhmump/backend-vet-clinic/internal/delivery/http"
	"github.com/mcmhmump/backend-vet-clinic/internal/repository"
	"github.com/mcmhmump/backend-vet-clinic/internal/usecase"
	"github.com/mcmhmump/backend-vet-clinic/pkg/logger"
	"go.uber.org/zap"
)

// @title Vet Clinic Proxy API
// @version 1.0
// @description API прокси-сервиса для управления IP-доступом, rate limiting, мониторингом и кэшированием.
// @host localhost:8000
// @BasePath /
func main() {
	appLogger := logger.New()
	defer appLogger.Sync()

	db, err := repository.NewDB()
	if err != nil {
		appLogger.Fatal("failed to connect database", zap.Error(err))
	}

	ipRuleRepo := repository.NewIPRuleRepository(db)
	accessLogRepo := repository.NewAccessLogRepository(db)

	ipRuleUsecase := usecase.NewIPRuleUsecase(ipRuleRepo)
	metricsService := usecase.NewMetricsService()
	rateLimiter := usecase.NewRateLimiter(5, 10*time.Second)
	cacheService := usecase.NewCacheService()

	ipRuleHandler := deliveryHttp.NewIPRuleHandler(ipRuleUsecase, appLogger, accessLogRepo, metricsService)
	monitoringHandler := deliveryHttp.NewMonitoringHandler(metricsService)
	cacheHandler := deliveryHttp.NewCacheHandler(cacheService)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /example/ip_access/allowlists", ipRuleHandler.GetAll)
	mux.HandleFunc("POST /example/ip_access/allowlists", ipRuleHandler.Create)
	mux.HandleFunc("DELETE /example/ip_access/allowlists/{id}", ipRuleHandler.Delete)
	mux.HandleFunc("GET /example/ip_access/check", ipRuleHandler.CheckIP)
	mux.HandleFunc("GET /example/metrics", monitoringHandler.GetMetrics)
	mux.HandleFunc("DELETE /example/cache", cacheHandler.InvalidateCache)

	mux.HandleFunc("GET /example/cache/demo", deliveryHttp.CacheMiddleware(cacheService, 30*time.Second, func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":   "slow data generated",
			"source":    "server",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	}))

	handlerWithRateLimit := usecase.RateLimitMiddleware(rateLimiter, mux)

	appLogger.Info("proxy api started", zap.String("port", ":8000"))

	if err := http.ListenAndServe(":8000", handlerWithRateLimit); err != nil {
		appLogger.Fatal("server failed", zap.Error(err))
	}
}
