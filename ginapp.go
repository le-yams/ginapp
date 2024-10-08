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

type App struct {
	logger    *zap.SugaredLogger
	config    Config
	ginEngine *gin.Engine
}

type Builder struct {
	config             Config
	configureGinEngine func(*gin.Engine, *zap.SugaredLogger) error
	configureMetrics   func(*ginmetrics.Monitor) error
}

func WithConfiguration(config Config) *Builder {
	return &Builder{
		config: config,
	}
}

func (builder *Builder) ConfigureGinEngine(configure func(*gin.Engine, *zap.SugaredLogger) error) *Builder {
	builder.configureGinEngine = configure
	return builder
}

func (builder *Builder) ConfigureMetrics(configure func(*ginmetrics.Monitor) error) *Builder {
	builder.configureMetrics = configure
	return builder
}

func (builder *Builder) Build() (*App, error) {
	config := builder.config
	logger, err := setupLogger(config)
	if err != nil {
		return nil, err
	}

	ginEngine, err := builder.setupGinEngine(config, logger)
	if err != nil {
		return nil, err
	}

	return &App{
		config:    config,
		logger:    logger,
		ginEngine: ginEngine,
	}, nil
}

func (builder *Builder) setupGinEngine(config Config, logger *zap.SugaredLogger) (*gin.Engine, error) {
	serverConfig := config.GetServerConfig()
	if serverConfig == nil {
		serverConfig = defaultServerConfig()
	}

	if serverConfig.Mode != "" {
		gin.SetMode(serverConfig.Mode)
	}

	ginEngine := gin.New()
	ginEngine.Use(ginzap.GinzapWithConfig(logger.Desugar(), &ginzap.Config{
		TimeFormat: time.RFC3339,
		UTC:        true,
		SkipPaths: []string{
			serverConfig.GetHealthCheckPath(),
			serverConfig.GetMetrics().GetPath(),
		},
	}))
	ginEngine.Use(ginzap.RecoveryWithZap(logger.Desugar(), true))
	ginEngine.Use(
		func(context *gin.Context) {
			context.Set("logger", logger.With("request_id", uuid.New().String()))
			context.Set("config", config)
			context.Next()
		},
	)

	setupHealthcheck(ginEngine, serverConfig)

	err := setupMetrics(ginEngine, serverConfig.GetMetrics(), builder.configureMetrics)
	if err != nil {
		return nil, fmt.Errorf("error configuring metrics: %w", err)
	}

	if builder.configureGinEngine != nil {
		err = builder.configureGinEngine(ginEngine, logger)
		if err != nil {
			return nil, fmt.Errorf("error configuring gin engine: %w", err)
		}
	}

	return ginEngine, nil
}

func (app *App) Start() error {
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

func (app *App) StartAsync() *http.Server {
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
