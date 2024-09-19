package ginapp

import "github.com/gin-gonic/gin"

const defaultMetricsPath = "/metrics"
const defaultHealthcheckPath = "/health"

type Config interface {
	GetServerConfig() *ServerConfig
	GetLogConfig() *LogConfig
}

type ServerConfig struct {
	Mode            string         `json:"mode,omitempty"`
	Port            int            `json:"port,omitempty"`
	HealthcheckPath string         `json:"healthcheck_path,omitempty"`
	MetricsEnabled  bool           `json:"enable_metrics,omitempty"`
	Metrics         *MetricsConfig `json:"metrics,omitempty"`
}

func defaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Mode:            gin.ReleaseMode,
		Port:            0,
		HealthcheckPath: defaultHealthcheckPath,
		Metrics:         defaultMetricsConfig(),
	}
}

func (serverConfig *ServerConfig) GetHealthCheckPath() string {
	if serverConfig.HealthcheckPath == "" {
		return defaultHealthcheckPath
	}
	return serverConfig.HealthcheckPath
}

func (serverConfig *ServerConfig) GetMetrics() *MetricsConfig {
	if serverConfig.Metrics == nil {
		serverConfig.Metrics = defaultMetricsConfig()
	}
	return serverConfig.Metrics
}

type MetricsConfig struct {
	Enabled bool   `json:"enabled,omitempty"`
	Path    string `json:"path,omitempty"`
}

func defaultMetricsConfig() *MetricsConfig {
	return &MetricsConfig{
		Path: defaultMetricsPath,
	}
}

func (metricsConfig *MetricsConfig) GetPath() string {
	if metricsConfig.Path == "" {
		return defaultMetricsPath
	}
	return metricsConfig.Path
}

type LogConfig struct {
	Level  LogLevel  `json:"level,omitempty"`
	Format LogFormat `json:"format,omitempty"`
}

func defaultLogConfig() *LogConfig {
	return &LogConfig{
		Level:  LogInfo,
		Format: LogConsole,
	}
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
