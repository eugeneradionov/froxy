package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DebugPreset = "debug"
	InfoPreset  = "info"
)

type Config struct {
	Level             zapcore.Level
	Development       bool
	DisableStacktrace bool
	EncoderConfig     zapcore.EncoderConfig
}

var DebugConfig = &Config{
	Level:             zapcore.DebugLevel,
	Development:       true,
	DisableStacktrace: false,
	EncoderConfig:     zap.NewDevelopmentEncoderConfig(),
}

var InfoConfig = &Config{
	Level:             zapcore.InfoLevel,
	Development:       false,
	DisableStacktrace: true,
	EncoderConfig:     zap.NewProductionEncoderConfig(),
}

func loadConfig(logPreset string) *Config {
	if logPreset == "debug" {
		return DebugConfig
	}

	return InfoConfig
}

func newConfig(cfg *Config) zap.Config {
	logConfig := defaultZapConfig()

	if cfg == nil {
		return logConfig
	}

	logConfig.Level = zap.NewAtomicLevelAt(cfg.Level)
	logConfig.Development = cfg.Development
	logConfig.DisableStacktrace = cfg.DisableStacktrace
	logConfig.EncoderConfig = cfg.EncoderConfig

	return logConfig
}

func defaultZapConfig() zap.Config {
	return zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:          "json",
		EncoderConfig:     zap.NewProductionEncoderConfig(),
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
		DisableStacktrace: false,
	}
}
