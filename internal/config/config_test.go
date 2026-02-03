package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultYAML(t *testing.T) {
	assert.Contains(t, DefaultYAML, "decisions_dir")
}
