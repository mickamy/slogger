package slogger_test

import (
	"io"
	"os"
	"testing"

	"github.com/mickamy/slogger"
)

func TestSlogger(t *testing.T) {
	t.Parallel()

	t.Run("DefaultConfig", func(t *testing.T) {
		t.Logf("Running test %s", t.Name())
		slogger.DebugCtx(t.Context(), "debug")
		slogger.InfoCtx(t.Context(), "info")
		slogger.WarnCtx(t.Context(), "warn")
		slogger.ErrorCtx(t.Context(), "error")
		t.Logf("Test %s completed", t.Name())
	})

	t.Run("CustomConfig", func(t *testing.T) {
		t.Logf("Running test %s", t.Name())
		slogger.SetConfig(slogger.Config{
			Level:          slogger.LevelDebug,
			AddSource:      true,
			TrimPathPrefix: "/Users/mickamy/ghq/github.com/mickamy/slogger",
			Outputs: []io.Writer{
				os.Stdout,
			},
			Handler: slogger.DefaultHandler(slogger.LevelDebug, os.Stdout),
		})

		slogger.DebugCtx(t.Context(), "debug with custom config")
		slogger.InfoCtx(t.Context(), "info with custom config")
		slogger.WarnCtx(t.Context(), "warn with custom config")
		slogger.ErrorCtx(t.Context(), "error with custom config")

		t.Logf("Test %s completed", t.Name())
	})
}
