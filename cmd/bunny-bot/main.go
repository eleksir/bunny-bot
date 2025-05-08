package main

import (
	"bunny-bot/internal/moon"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	var err error

	// Main config.
	executablePath, err := os.Executable()

	if err != nil {
		panic(fmt.Sprintf("Unable to get current executable path: %s", err))
	}

	configJSONPath := fmt.Sprintf("%s/data/config.json", filepath.Dir(executablePath))

	moon.Cfg, err = moon.ParseConfig(configJSONPath)

	if err != nil {
		panic(err)
	}

	// Config for Zap logger facility.
	var logCfg zap.Config = zap.NewDevelopmentConfig()

	// Define loglevel.
	var l zapcore.Level

	if err := l.UnmarshalText([]byte(moon.Cfg.LogLevel)); err != nil {
		panic("Unknown LogLevel")
	}

	logCfg.Level = zap.NewAtomicLevelAt(l)

	// Define output file if set in config.
	if moon.Cfg.LogFile != "" {
		logCfg.OutputPaths[0] = moon.Cfg.LogFile
	}

	// We do not want to see caller in logs.
	logCfg.DisableCaller = true

	// We want human-readable timestamp in logs.
	logCfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)

	logger := zap.Must(logCfg.Build())
	moon.Log = logger.Sugar()

	mystdlogger := zap.RedirectStdLog(logger)
	defer mystdlogger()

	moon.Log.Info("Config read looks good")

	// Самое время поставить траппер сигналов.
	signal.Notify(moon.SigChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go moon.SigHandler()

	moon.Log.Info("Singnal handler installed")

	moon.Telegram(moon.Cfg)
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
