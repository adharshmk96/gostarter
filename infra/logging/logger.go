package logging

import (
	"log/slog"
	"os"
)

func NewFileLogger(path string) *slog.Logger {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return slog.New(slog.NewJSONHandler(file, nil))
}
