package main

import (
	"fmt"
	"os"

	"github.com/br9k777/calendar/pkg/config"
	"github.com/br9k777/calendar/pkg/daemon"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	logger, err := config.GetStandartLogger("production")
	if err != nil {
		fmt.Fprintf(os.Stdout, "Can't create logger %s", err)
		os.Exit(1)
	}
	zap.ReplaceGlobals(logger)

	config.InitDefaultConfigSettings()
	cfg, err := config.GetConfig()
	if err != nil {
		zap.S().Fatal(err)
	}
	if err = viper.WriteConfigAs(".config.yaml"); err != nil {
		zap.S().Error(err)
	}

	var newLogger *zap.Logger
	if cfg.LogConfigurationFile != "" {
		if newLogger, err = config.GetLoggerConfigFromFile(cfg.LogConfigurationFile); err != nil {
			zap.S().Error(err)
		}
	}
	if newLogger == nil {
		if newLogger, err = config.GetStandartLogger(cfg.LogStandartChois); err != nil {
			zap.S().Fatal(err)
		}
	}
	if err = daemon.Run(cfg, newLogger); err != nil {
		newLogger.Error("Daemon run error", zap.Error(err))
	}

}
