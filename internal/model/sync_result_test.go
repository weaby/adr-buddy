package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSyncResult_MarshalJSON(t *testing.T) {
	result := &SyncResult{
		ChangesDetected: true,
		Files: FileChanges{
			Created:  []string{"decisions/infrastructure/adr-5.md"},
			Modified: []string{"decisions/adr-1.md"},
			Deleted:  []string{},
		},
		ADRs: []ADRChange{
			{
				ID:       "adr-5",
				Name:     "Using Redis for caching",
				Action:   "create",
				FilePath: "decisions/infrastructure/adr-5.md",
			},
		},
	}

	data, err := json.Marshal(result)
	require.NoError(t, err)

	var decoded SyncResult
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)
	assert.True(t, decoded.ChangesDetected)
	assert.Equal(t, 1, len(decoded.Files.Created))
	assert.Equal(t, 1, len(decoded.ADRs))
}
