package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

func SetupLogging(level string) error {
	slogLevel, err := parseLogLevel(level)
	if err != nil {
		return err
	}

	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slogLevel,
	})
	slog.SetDefault(slog.New(handler))

	return nil
}

func parseLogLevel(level string) (slog.Level, error) {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return 0, fmt.Errorf("invalid log-level: %s", level)
	}
}
