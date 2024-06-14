package loggerutilities

import (
	"log/slog"
	"os"

	"github.com/Originate/go-utilities/configutilities"
)

func SetupLogger(cfg configutilities.SlogConfiguration) {
	var logLevel = new(slog.LevelVar)
	if err := logLevel.UnmarshalText([]byte(cfg.Level)); err != nil {
		slog.Warn("Unrecognized slog level, using default", "level", cfg.Level, "default", logLevel.Level())
	}

	logger := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: logLevel})
	slog.SetDefault(slog.New(logger))

	slog.Info("Setting log", "level", logLevel)
	logLevel.Set(logLevel.Level())
}
