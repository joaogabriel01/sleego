package logger

import "io"

// LoggerMock is a mock implementation of the Logger interface for testing purposes.'
type LoggerMock struct {
}

func NewLoggerMock() *LoggerMock {
	return &LoggerMock{}
}

func (l *LoggerMock) Info(msg string) {}

func (l *LoggerMock) Debug(msg string) {}

func (l *LoggerMock) Error(msg string) {}

func (l *LoggerMock) WithField(key string, value interface{}) Logger {
	return l
}

func (l *LoggerMock) Out() io.Writer {
	return nil
}
