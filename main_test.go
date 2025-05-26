package slogger_test

import (
	"testing"

	"github.com/mickamy/slogger"
)

func TestSlogger(t *testing.T) {
	t.Parallel()

	t.Logf("Running test %s", t.Name())
	slogger.DebugCtx(t.Context(), "debug")
	slogger.InfoCtx(t.Context(), "info")
	slogger.WarnCtx(t.Context(), "warn")
	slogger.ErrorCtx(t.Context(), "error")
	t.Logf("Test %s completed", t.Name())
}
