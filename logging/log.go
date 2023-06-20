package logging

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/exp/slog"
)

func SetupLoging(name string) {
	// setup logging
	f, err := os.Create(os.TempDir() + "/" + name + ".log")
	if err != nil {
		log.Println("log file creation", err)
	}
	//defer f.Close() -- don't close file here
	logLevel := &slog.LevelVar{}
	replace := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.SourceKey {
			a.Value = slog.StringValue(filepath.Base(a.Value.String()))
		}
		return a
	}
	logger := slog.New(slog.NewJSONHandler(io.MultiWriter(os.Stderr, f), &slog.HandlerOptions{AddSource: true, ReplaceAttr: replace, Level: logLevel}))
	logger2 := logger.With("name", name)
	slog.SetDefault(logger2)
}
