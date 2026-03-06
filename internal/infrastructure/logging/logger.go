package logging

import (
	"os"
	"team_service/internal/infrastructure/share/settings"

	"github.com/thanvuc/go-core-lib/log"
)

func NewLogger(config *settings.Log) log.Logger {
	env := os.Getenv("GO_ENV")
	return log.NewLogger(log.Config{
		Env:   env,
		Level: config.Level,
	})
}

func NewLoggerV2() (log.LoggerV2, error) {
	env := os.Getenv("GO_ENV")
	println("Logger creating with env:", env)
	return log.NewLoggerZapV2(env)
}
