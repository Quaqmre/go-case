package util

import (
	"os"

	"github.com/go-kit/log/level"
)

func GetLogLevelFromEnv() level.Option {
	logEnv, ok := os.LookupEnv("LOG_LEVEL")
	// LOG_LEVEL is not set, let's default to debug
	if !ok {
		return level.AllowAll()
	}
	levelValue := level.ParseDefault(logEnv, level.DebugValue())
	switch levelValue {
	case level.ErrorValue():
		return level.AllowError()
	case level.DebugValue():
		return level.AllowDebug()
	case level.InfoValue():
		return level.AllowInfo()
	default:
		return level.AllowAll()
	}
}
