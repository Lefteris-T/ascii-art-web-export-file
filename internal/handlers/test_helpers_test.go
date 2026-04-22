package handlers

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestMain(m *testing.M) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		os.Exit(1)
	}

	repoRoot := filepath.Dir(filepath.Dir(filepath.Dir(thisFile)))
	if err := os.Chdir(repoRoot); err != nil {
		os.Exit(1)
	}

	os.Exit(m.Run())
}
