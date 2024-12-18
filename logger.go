package main

import (
	"os"
	"path"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

func setLogLevel() {
	switch config.LogLevel {
	case "DEBUG":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

func getLogger() zerolog.Logger {
	wd, err := os.Getwd()
	processErrorWithoutLogger(err, "Error setting up logger")
	logFile := &lumberjack.Logger{
		Filename:   path.Join(wd, "excel-automations.log"),
		MaxSize:    500,
		MaxBackups: 2,
	}
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339Nano}
	multiWriter := zerolog.MultiLevelWriter(consoleWriter, logFile)
	setLogLevel()
	return zerolog.New(multiWriter).With().Timestamp().Logger()
}
