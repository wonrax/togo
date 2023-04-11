package togo

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func initLogger(env string) error {
	var logConfig = zap.NewDevelopmentConfig()
	logConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	if env == "production" || env == "test" {
		logConfig = zap.NewProductionConfig()
		logConfig.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	}

	var err error
	Log, err = logConfig.Build()

	return err
}
