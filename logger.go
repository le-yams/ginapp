package ginapp

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func setupLogger(config Config) (*zap.SugaredLogger, error) {
	logConfig := config.GetLogConfig()

	if logConfig == nil {
		logConfig = defaultLogConfig()
	}
	if logConfig.Level == LogNone {
		return zap.NewNop().Sugar(), nil
	}

	loggerConfig := zap.NewProductionConfig()
	if config.GetServerConfig().Mode == gin.ReleaseMode {
		loggerConfig = zap.NewDevelopmentConfig()
	}
	loggerConfig.Encoding = string(logConfig.Format)
	loggerConfig.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder

	switch logConfig.Level {
	case LogDebug:
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case LogInfo:
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case LogWarn:
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case LogError:
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	}

	logger, err := loggerConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("error building logger: %w", err)
	}

	return logger.Sugar(), nil
}
