package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog"
)

var globalLogger Logger

type Logger interface {
	Info(msg string)
	Debug(msg string)
	Error(msg string)
	WithField(key string, value interface{}) Logger
	Out() io.Writer
}

type ZeroLogger struct {
	logger zerolog.Logger
	output io.Writer
}

func Get(logLevel ...string) (Logger, error) {
	if globalLogger != nil {
		return globalLogger, nil
	}

	logLevelStr := "info"
	if len(logLevel) > 0 {
		logLevelStr = logLevel[0]
	}

	initLogger(logLevelStr)
	if globalLogger != nil {
		return globalLogger, nil
	}

	return nil, fmt.Errorf("global logger is not initialized")
}

func initLogger(logLevel string) {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05"}

	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}

	logger := zerolog.New(output).Level(level).With().Timestamp().Logger()

	globalLogger = &ZeroLogger{
		logger: logger,
		output: output,
	}
}

func (zl *ZeroLogger) Info(msg string) {
	zl.logger.Info().Msg(msg)
}

func (zl *ZeroLogger) Debug(msg string) {
	zl.logger.Debug().Msg(msg)
}

func (zl *ZeroLogger) Error(msg string) {
	zl.logger.Error().Msg(msg)
}

func (zl *ZeroLogger) WithField(key string, value interface{}) Logger {
	newLogger := zl.logger.With().Interface(key, value).Logger()
	return &ZeroLogger{
		logger: newLogger,
		output: zl.output,
	}
}

func (zl *ZeroLogger) Out() io.Writer {
	return zl.output
}
