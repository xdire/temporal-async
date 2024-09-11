package util

import (
	"github.com/rs/zerolog"
	"os"
)

func NewZeroLogForName(name, id, level string) zerolog.Logger {
	zLevel := zerolog.ErrorLevel
	if len(level) > 0 {
		newLevel, err := zerolog.ParseLevel(level)
		if err == nil {
			zLevel = newLevel
		}
	}
	return zerolog.New(os.Stdout).
		Level(zLevel).With().Timestamp().
		Caller().Str(name, id).Logger()
}
