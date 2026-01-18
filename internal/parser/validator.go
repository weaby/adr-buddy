package parser

import (
	"fmt"
)

var validStatuses = map[string]bool{
	"proposed":   true,
	"accepted":   true,
	"rejected":   true,
	"deprecated": true,
	"superseded": true,
}

// ValidateStatus checks if a status value is valid
func ValidateStatus(status string, strict bool) error {
	if status == "" {
		return nil // Empty is OK, will default to "proposed"
	}

	if !validStatuses[status] {
		if strict {
			return fmt.Errorf("invalid status %q: must be one of: proposed, accepted, rejected, deprecated, superseded", status)
		}
		// Just a warning in non-strict mode
		fmt.Printf("WARNING: Unknown status %q\n", status)
	}

	return nil
}
