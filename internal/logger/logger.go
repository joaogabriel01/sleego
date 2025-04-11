package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Init(logLevel string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05"}

	log.Logger = zerolog.New(output).Level(level).With().Timestamp().Logger()
}

func Get() zerolog.Logger {
	return log.Logger
}
