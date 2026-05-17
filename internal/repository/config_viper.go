package repository

import (
	"github.com/fsnotify/fsnotify"
	"github.com/mcmhmump/backend-vet-clinic/internal/domain"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func LoadConfig(logger *zap.Logger, onChange func(cfg *domain.AppConfig)) (*domain.AppConfig, error) {
	viper.SetConfigFile("config.yaml")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config domain.AppConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		logger.Info("config file changed", zap.String("file", e.Name))

		var newConfig domain.AppConfig
		if err := viper.Unmarshal(&newConfig); err != nil {
			logger.Error("failed to unmarshal updated config", zap.Error(err))
			return
		}

		if onChange != nil {
			onChange(&newConfig)
		}
	})

	viper.WatchConfig()

	return &config, nil
}
