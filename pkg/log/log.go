// Copyright (c) 2019 IoTeX Foundation
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

// A warrper for Zerolog (https://github.com/rs/zerolog)
//
// Package log provides a global logger for zerolog.
// derived from https://github.com/rs/zerolog/blob/master/log/log.go
// putting here to get a better integration

package log

import (
	"log"
	"net/http"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// GlobalConfig defines the global logger configurations.
type GlobalConfig struct {
	Zap                *zap.Config `json:"zap" yaml:"zap"`
	StderrRedirectFile *string     `json:"stderrRedirectFile" yaml:"stderrRedirectFile"`
	RedirectStdLog     bool        `json:"stdLogRedirect" yaml:"stdLogRedirect"`
}

var (
	_globalCfg        GlobalConfig
	_logMu            sync.RWMutex
	_logServeMux      = http.NewServeMux()
	_subLoggers       map[string]*zap.Logger
	_globalLoggerName = "global"
)

func init() {
	zapCfg := zap.NewDevelopmentConfig()
	zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	zapCfg.Level.SetLevel(zap.InfoLevel)
	l, err := zapCfg.Build()
	if err != nil {
		log.Println("Failed to init zap global logger, no zap log will be shown till zap is properly initialized: ", err)
		return
	}
	_logMu.Lock()
	_globalCfg.Zap = &zapCfg
	_logMu.Unlock()
	zap.ReplaceGlobals(l)
}

// L wraps zap.L().
func L() *zap.Logger { return zap.L() }

// S wraps zap.S().
func S() *zap.SugaredLogger { return zap.S() }

// Logger returns logger of the given name
func Logger(name string) *zap.Logger {
	logger, ok := _subLoggers[name]
	if !ok {
		return L()
	}
	return logger
}

// InitLoggers initializes the global logger and other sub loggers.
func InitLoggers(globalCfg GlobalConfig) error {
	cfg := globalCfg
	if cfg.Zap == nil {
		zapCfg := zap.NewProductionConfig()
		cfg.Zap = &zapCfg
	} else {
		cfg.Zap.EncoderConfig = zap.NewProductionEncoderConfig()
	}

	if cfg.StderrRedirectFile != nil {
		stderrF, err := os.OpenFile(*cfg.StderrRedirectFile, os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		if err := redirectStderr(stderrF); err != nil {
			return err
		}
	}
	return nil
}
