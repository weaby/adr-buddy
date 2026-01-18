package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAggregate_SingleAnnotation(t *testing.T) {
	annotations := []*Annotation{
		{
			ID:       "adr-1",
			Name:     "Using Pino",
			Status:   "accepted",
			Category: "infrastructure",
			Context:  "Need fast logging",
			Location: SourceLocation{
				File: "src/logger.js",
				Line: 10,
			},
		},
	}

	adrs, err := Aggregate(annotations)
	assert.NoError(t, err)
	assert.Len(t, adrs, 1)

	adr := adrs[0]
	assert.Equal(t, "adr-1", adr.ID)
	assert.Equal(t, "Using Pino", adr.Name)
	assert.Equal(t, "accepted", adr.Status)
	assert.Equal(t, "infrastructure", adr.Category)
	assert.Len(t, adr.Context, 1)
	assert.Equal(t, "Need fast logging", adr.Context[0])
	assert.Len(t, adr.Locations, 1)

	// Date should be today
	today := time.Now().Format("2006-01-02")
	assert.Equal(t, today, adr.Date)
}

func TestAggregate_MultipleAnnotations(t *testing.T) {
	annotations := []*Annotation{
		{
			ID:       "adr-1",
			Name:     "Using Pino",
			Context:  "Context from first location",
			Location: SourceLocation{File: "src/logger.js", Line: 10},
		},
		{
			ID:       "adr-1",
			Name:     "Using Pino",
			Decision: "Decision from second location",
			Location: SourceLocation{File: "src/app.js", Line: 5},
		},
	}

	adrs, err := Aggregate(annotations)
	assert.NoError(t, err)
	assert.Len(t, adrs, 1)

	adr := adrs[0]
	assert.Len(t, adr.Context, 1)
	assert.Equal(t, "Context from first location", adr.Context[0])
	assert.Len(t, adr.Decision, 1)
	assert.Equal(t, "Decision from second location", adr.Decision[0])
	assert.Len(t, adr.Locations, 2)
}

func TestAggregate_ConflictingNames(t *testing.T) {
	annotations := []*Annotation{
		{
			ID:       "adr-1",
			Name:     "Using Pino",
			Location: SourceLocation{File: "src/logger.js", Line: 10},
		},
		{
			ID:       "adr-1",
			Name:     "Different Name",
			Location: SourceLocation{File: "src/app.js", Line: 5},
		},
	}

	_, err := Aggregate(annotations)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "conflicting names")
}
