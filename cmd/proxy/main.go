package main

import (
	"net/http"

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

	ipRuleRepo := repository.NewIPRuleRepository(db)
	ipRuleUsecase := usecase.NewIPRuleUsecase(ipRuleRepo)
	ipRuleHandler := deliveryHttp.NewIPRuleHandler(ipRuleUsecase, appLogger)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /example/ip_access/allowlists", ipRuleHandler.GetAll)
	mux.HandleFunc("POST /example/ip_access/allowlists", ipRuleHandler.Create)
	mux.HandleFunc("DELETE /example/ip_access/allowlists/{id}", ipRuleHandler.Delete)
	mux.HandleFunc("GET /example/ip_access/check", ipRuleHandler.CheckIP)

	appLogger.Info("management api started", zap.String("port", ":8000"))

	if err := http.ListenAndServe(":8000", mux); err != nil {
		appLogger.Fatal("server failed", zap.Error(err))
	}
}
