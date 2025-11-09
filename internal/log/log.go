package log

import (
	"log/slog"
	"os"
)

func NewLogger(cfg interface{}) *slog.Logger {
	// TODO: depending on config, change should be possible to io.Writer.
	handler := slog.NewJSONHandler(os.Stdout, nil)

	logger := slog.New(handler)

	return logger
}
