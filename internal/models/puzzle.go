package models

import (
	"time"
)

// Cell represents a single cell in the Sudoku grid
type Cell struct {
	Value  int    `json:"value"`  // 1-9, 0 means unset
	Notes  []int  `json:"notes"`  // Array of 1-9 values
	Status string `json:"status"` // s=system, u=user, w=wrong, c=correct
}

// Puzzle represents a Sudoku puzzle with metadata
type Puzzle struct {
	UUID       string          `json:"uuid"`
	CreatedAt  string          `json:"createdAt"`
	Cells      map[string]Cell `json:"cells"` // Position (01-81) as key
	Difficulty int             `json:"difficulty"`
}

// GetTimeString returns the current time in RFC3339 format
func GetTimeString() string {
	return time.Now().Format(time.RFC3339)
}
