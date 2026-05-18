package main

import (
	"net/http"
	"time"

	deliveryHttp "github.com/mcmhmump/backend-vet-clinic/internal/delivery/http"
	"github.com/mcmhmump/backend-vet-clinic/internal/repository"
	"github.com/mcmhmump/backend-vet-clinic/internal/usecase"
	"github.com/mcmhmump/backend-vet-clinic/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	appLogger := logger.New()
	defer appLogger.Sync()

	db, err := repository.NewDB()
	if err != nil {
		appLogger.Fatal("failed to connect database", zap.Error(err))
	}

	// Инициализация репозиториев
	ipRuleRepo := repository.NewIPRuleRepository(db)
	accessLogRepo := repository.NewAccessLogRepository(db)

	// Инициализация бизнес-логики (Usecase / Services)
	ipRuleUsecase := usecase.NewIPRuleUsecase(ipRuleRepo)
	metricsService := usecase.NewMetricsService()
	rateLimiter := usecase.NewRateLimiter(5, 10*time.Second)

	// Инициализация HTTP обработчиков
	ipRuleHandler := deliveryHttp.NewIPRuleHandler(ipRuleUsecase, appLogger, accessLogRepo, metricsService)
	monitoringHandler := deliveryHttp.NewMonitoringHandler(metricsService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /example/ip_access/allowlists", ipRuleHandler.GetAll)
	mux.HandleFunc("POST /example/ip_access/allowlists", ipRuleHandler.Create)
	mux.HandleFunc("DELETE /example/ip_access/allowlists/{id}", ipRuleHandler.Delete)
	mux.HandleFunc("GET /example/ip_access/check", ipRuleHandler.CheckIP)

	// НОВЫЙ ЭНДПОИНТ ДЛЯ МЕТРИК
	mux.HandleFunc("GET /example/metrics", monitoringHandler.GetMetrics)

	// Оборачиваем роутер в Rate Limiting Middleware
	handlerWithRateLimit := usecase.RateLimitMiddleware(rateLimiter, mux)

	appLogger.Info("proxy and monitoring api started", zap.String("port", ":8000"))

	if err := http.ListenAndServe(":8000", handlerWithRateLimit); err != nil {
		appLogger.Fatal("server failed", zap.Error(err))
	}
}
