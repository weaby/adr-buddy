package main

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExitCodes(t *testing.T) {
	// Build the binary
	cmd := exec.Command("go", "build", "-o", "/tmp/adr-buddy-test", ".")
	err := cmd.Run()
	assert.NoError(t, err)
	defer os.Remove("/tmp/adr-buddy-test")

	// Test successful check - exit 0
	tmpDir := t.TempDir()
	os.WriteFile(tmpDir+"/.adr-buddy/config.yml", []byte("scan_paths: [.]"), 0644)

	cmd = exec.Command("/tmp/adr-buddy-test", "check")
	cmd.Dir = tmpDir
	err = cmd.Run()
	assert.NoError(t, err) // Exit code 0

	// Test failed check - exit 1
	// (would need invalid annotations to test)
}
