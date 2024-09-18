package appengine

import (
	"errors"
	"fmt"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"log"
	"net"
	"net/http"
	"time"
)

type GinAppConfig interface {
	GetServerConfig() *ServerConfig
	GetLogConfig() *LogConfig
}

type ServerConfig struct {
	Port int `json:"port,omitempty"`
}

type LogConfig struct {
	Level  LogLevel  `json:"level,omitempty"`
	Format LogFormat `json:"format,omitempty"`
}

type LogLevel string

const (
	LogNone  LogLevel = "none"
	LogDebug LogLevel = "debug"
	LogInfo  LogLevel = "info"
	LogWarn  LogLevel = "warn"
	LogError LogLevel = "error"
)

type LogFormat string

const (
	LogJson    LogFormat = "console"
	LogConsole LogFormat = "json"
)

type GinApp struct {
	logger    *zap.SugaredLogger
	config    GinAppConfig
	ginEngine *gin.Engine
}

func New(configuration GinAppConfig, configure func(*gin.Engine, *zap.SugaredLogger) error) (*GinApp, error) {
	logger, err := setupLogger(configuration)
	if err != nil {
		return nil, err
	}

	if configuration.GetLogConfig().Level == LogDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	ginEngine := gin.New()

	ginEngine.Use(ginzap.Ginzap(logger.Desugar(), time.RFC3339, true))
	ginEngine.Use(ginzap.RecoveryWithZap(logger.Desugar(), true))
	ginEngine.Use(
		func(context *gin.Context) {
			context.Set("logger", logger.With("request_id", uuid.New().String()))
			context.Set("config", configuration)
			context.Next()
		},
	)

	err = configure(ginEngine, logger)
	if err != nil {
		return nil, err
	}

	return &GinApp{
		config:    configuration,
		logger:    logger,
		ginEngine: ginEngine,
	}, nil
}

func setupLogger(configuration GinAppConfig) (*zap.SugaredLogger, error) {
	logConfiguration := configuration.GetLogConfig()
	if logConfiguration.Level == LogNone {
		return zap.NewNop().Sugar(), nil
	}

	loggerConfig := zap.NewProductionConfig()
	loggerConfig.Encoding = string(logConfiguration.Format)

	switch logConfiguration.Level {
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
		return nil, err
	}

	return logger.Sugar(), nil
}

func (app *GinApp) Start() error {
	defer func() {
		if err := app.logger.Sync(); err != nil {
			app.logger.Warnw("cannot flush logger", "err", err)
		}
	}()
	port := app.config.GetServerConfig().Port
	address := fmt.Sprintf("localhost:%d", port)

	return app.ginEngine.Run(address)
}

func (app *GinApp) StartAsync() *http.Server {
	port := app.config.GetServerConfig().Port
	address := fmt.Sprintf("localhost:%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		app.logger.Fatalw("a fatal error occured when configuring HTTP listener", "err", err)
	}

	port = listener.Addr().(*net.TCPAddr).Port
	app.config.GetServerConfig().Port = port
	address = fmt.Sprintf("localhost:%d", port)

	app.logger.Infow("Configuring http listener", "addr", "localhost", "port", port)

	srv := &http.Server{
		Addr:    address,
		Handler: app.ginEngine,
	}

	go func() {
		defer func() {
			if err = app.logger.Sync(); err != nil {
				app.logger.Warnw("cannot flush logger", "err", err)
			}
		}()

		if err = srv.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	return srv
}
