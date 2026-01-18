package model

// FileChanges tracks what files will be created, modified, or deleted
type FileChanges struct {
	Created  []string `json:"created"`
	Modified []string `json:"modified"`
	Deleted  []string `json:"deleted"`
}

// ADRChange represents a single ADR that will be changed
type ADRChange struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Action   string `json:"action"` // "create", "update", "delete"
	FilePath string `json:"file_path"`
}

// SyncResult represents the output of sync --dry-run command
type SyncResult struct {
	ChangesDetected bool        `json:"changes_detected"`
	Files           FileChanges `json:"files"`
	ADRs            []ADRChange `json:"adrs"`
}
