package model

import (
	"fmt"
	"time"
)

// Aggregate combines annotations with the same ID into ADRs
func Aggregate(annotations []*Annotation) ([]*ADR, error) {
	adrMap := make(map[string]*ADR)

	for _, ann := range annotations {
		if err := ann.Validate(); err != nil {
			return nil, fmt.Errorf("%s: %w", ann.Location, err)
		}

		adr, exists := adrMap[ann.ID]
		if !exists {
			// Create new ADR
			adr = &ADR{
				ID:           ann.ID,
				Name:         ann.Name,
				Status:       ann.DefaultStatus(),
				Category:     ann.Category,
				Date:         time.Now().Format("2006-01-02"),
				Context:      []string{},
				Decision:     []string{},
				Alternatives: []string{},
				Consequences: []string{},
				Locations:    []SourceLocation{},
			}
			adrMap[ann.ID] = adr
		} else {
			// Validate consistency
			if adr.Name != ann.Name {
				return nil, fmt.Errorf("conflicting names for %s: %q at %s vs %q at %s",
					ann.ID, adr.Name, adr.Locations[0], ann.Name, ann.Location)
			}
			if adr.Category != ann.Category {
				return nil, fmt.Errorf("conflicting categories for %s: %q at %s vs %q at %s",
					ann.ID, adr.Category, adr.Locations[0], ann.Category, ann.Location)
			}
		}

		// Append content fields
		if ann.Context != "" {
			adr.Context = append(adr.Context, ann.Context)
		}
		if ann.Decision != "" {
			adr.Decision = append(adr.Decision, ann.Decision)
		}
		if ann.Alternatives != "" {
			adr.Alternatives = append(adr.Alternatives, ann.Alternatives)
		}
		if ann.Consequences != "" {
			adr.Consequences = append(adr.Consequences, ann.Consequences)
		}

		// Add location
		adr.Locations = append(adr.Locations, ann.Location)
	}

	// Convert map to slice
	adrs := make([]*ADR, 0, len(adrMap))
	for _, adr := range adrMap {
		adrs = append(adrs, adr)
	}

	return adrs, nil
}
