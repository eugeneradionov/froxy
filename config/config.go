package config

import (
	"errors"
	"fmt"

	"github.com/jessevdk/go-flags"
)

var config Config

func Get() *Config {
	return &config
}

func Load() error {
	p := flags.NewParser(&config, flags.PrintErrors|flags.PassDoubleDash|flags.HelpFlag)
	p.SubcommandsOptional = true

	if _, err := p.Parse(); err != nil {
		if !errors.Is(err.(*flags.Error).Type, flags.ErrHelp) { // nolint:errorlint
			return fmt.Errorf("cli error: %w", err)
		}
	}

	return nil
}

type Config struct {
	Logger     Logger     `group:"logger" namespace:"logger" env-namespace:"LOGGER"`
	Proxy      Proxy      `group:"proxy" namespace:"proxy" env-namespace:"PROXY"`
	HTTPServer HTTPServer `group:"http-server" namespace:"http-server" env-namespace:"HTTP_SERVER"`
}

type Logger struct {
	Preset string `long:"preset" env:"PRESET" default:"info" description:"logger preset (debug|info)"`
}

type HTTPServer struct {
	ListenURL string `long:"listen-url" env:"LISTEN_URL" default:"localhost:8080" description:"server listen url"`
}

type Proxy struct {
	FileMaxSizeMB int64 `long:"file-max-size" env:"FILE_MAX_SIZE" default:"5" description:"max file size for upload in MB"`
}
