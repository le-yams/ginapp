package ginapp

import (
	"errors"
	//"errors"
	"fmt"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/penglongli/gin-metrics/ginmetrics"
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
	Port            int    `json:"port,omitempty"`
	HealthcheckPath string `json:"healthcheck_path,omitempty"`
	MetricsEnabled  bool   `json:"enable_metrics,omitempty"`
}

const DefaultHealthcheckPath = "/health"

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
	LogJson    LogFormat = "json"
	LogConsole LogFormat = "console"
)

type GinApp struct {
	logger    *zap.SugaredLogger
	config    GinAppConfig
	ginEngine *gin.Engine
}

type Setups interface {
	ConfigureGinEngine(*gin.Engine, *zap.SugaredLogger) error
	ConfigureMetrics(*ginmetrics.Monitor) error
}

func New(configuration GinAppConfig, setups Setups) (*GinApp, error) {
	logConfiguration := configuration.GetLogConfig()
	if logConfiguration == nil {
		logConfiguration = &LogConfig{
			Level:  LogInfo,
			Format: LogConsole,
		}
	}
	logger, err := setupLogger(logConfiguration)
	if err != nil {
		return nil, err
	}

	if logConfiguration.Level == LogDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	ginEngine, _ := setupGinEngine(configuration, setups, logger)

	return &GinApp{
		config:    configuration,
		logger:    logger,
		ginEngine: ginEngine,
	}, nil
}

func setupLogger(configuration *LogConfig) (*zap.SugaredLogger, error) {
	if configuration.Level == LogNone {
		return zap.NewNop().Sugar(), nil
	}

	loggerConfig := zap.NewProductionConfig()
	loggerConfig.Encoding = string(configuration.Format)

	switch configuration.Level {
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

func setupGinEngine(configuration GinAppConfig, setups Setups, logger *zap.SugaredLogger) (*gin.Engine, error) {
	ginEngine := gin.New()

	healthcheckPath := configuration.GetServerConfig().HealthcheckPath
	if healthcheckPath == "" {
		healthcheckPath = DefaultHealthcheckPath
	}
	ginEngine.Use(ginzap.GinzapWithConfig(logger.Desugar(), &ginzap.Config{
		TimeFormat: time.RFC3339,
		UTC:        true,
		SkipPaths:  []string{healthcheckPath},
	}))
	ginEngine.Use(ginzap.RecoveryWithZap(logger.Desugar(), true))
	ginEngine.Use(
		func(context *gin.Context) {
			context.Set("logger", logger.With("request_id", uuid.New().String()))
			context.Set("config", configuration)
			context.Next()
		},
	)

	ginEngine.GET(healthcheckPath, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	if setups != nil {
		if configuration.GetServerConfig().MetricsEnabled {
			err := setupMetrics(ginEngine, setups.ConfigureMetrics)
			if err != nil {
				return nil, err
			}
		}

		err := setups.ConfigureGinEngine(ginEngine, logger)
		if err != nil {
			return nil, fmt.Errorf("error configuring gin engine: %w", err)
		}
	}

	return ginEngine, nil
}

func setupMetrics(ginEngine *gin.Engine, configure func(*ginmetrics.Monitor) error) error {
	monitor := ginmetrics.GetMonitor()
	monitor.SetMetricPath("/metrics")
	err := configure(monitor)
	if err != nil {
		return err
	}
	monitor.Use(ginEngine)
	return nil
}

func (app *GinApp) Start() error {
	defer func() {
		if err := app.logger.Sync(); err != nil {
			app.logger.Warnw("cannot flush logger", "err", err)
		}
	}()
	port := app.config.GetServerConfig().Port
	address := fmt.Sprintf("localhost:%d", port)

	err := app.ginEngine.Run(address)
	if err != nil {
		return fmt.Errorf("application error: %w", err)
	}

	return nil
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
		Addr:              address,
		Handler:           app.ginEngine,
		ReadHeaderTimeout: 10 * time.Second,
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
