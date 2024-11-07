package sl

import (
	"log/slog"
	"os"
)

func InitLogger() *slog.Logger {
	var log *slog.Logger

	log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	log.Info("Logger OK")

	return log
}
