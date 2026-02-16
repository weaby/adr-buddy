package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExitCodes(t *testing.T) {
	// Build the binary
	tmpBin := filepath.Join(t.TempDir(), "adr-buddy")
	build := exec.Command("go", "build", "-o", tmpBin, ".")
	build.Dir = "."
	out, err := build.CombinedOutput()
	require.NoError(t, err, "build failed: %s", string(out))

	t.Run("check succeeds with no ADRs", func(t *testing.T) {
		tmpDir := t.TempDir()
		cmd := exec.Command(tmpBin, "check")
		cmd.Dir = tmpDir
		err := cmd.Run()
		require.NoError(t, err)
	})

	t.Run("list succeeds with no ADRs", func(t *testing.T) {
		tmpDir := t.TempDir()
		cmd := exec.Command(tmpBin, "list")
		cmd.Dir = tmpDir
		err := cmd.Run()
		require.NoError(t, err)
	})

	t.Run("new requires name argument", func(t *testing.T) {
		tmpDir := t.TempDir()
		cmd := exec.Command(tmpBin, "new")
		cmd.Dir = tmpDir
		err := cmd.Run()
		require.Error(t, err) // Should fail without args
	})
}

func TestExitCodes_Init(t *testing.T) {
	// Build the binary
	tmpBin := filepath.Join(t.TempDir(), "adr-buddy")
	build := exec.Command("go", "build", "-o", tmpBin, ".")
	build.Dir = "."
	out, err := build.CombinedOutput()
	require.NoError(t, err, "build failed: %s", string(out))

	t.Run("init with skip creates decisions dir", func(t *testing.T) {
		tmpDir := t.TempDir()
		cmd := exec.Command(tmpBin, "init", "--claude-skill", "skip")
		cmd.Dir = tmpDir
		err := cmd.Run()
		require.NoError(t, err)

		decisionsDir := filepath.Join(tmpDir, ".claude", "rules", "decisions")
		_, err = os.Stat(decisionsDir)
		require.NoError(t, err)
	})
}
