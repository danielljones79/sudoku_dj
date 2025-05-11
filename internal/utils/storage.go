package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/danjones/sudoku_dj/internal/models"
)

// SavePuzzle saves a puzzle to disk
func SavePuzzle(puzzle models.Puzzle) (models.Puzzle, error) {
	data, err := json.MarshalIndent(puzzle, "", "  ")
	if err != nil {
		Log(LogLevelError, "Error marshaling puzzle: %v", err)
		return puzzle, fmt.Errorf("error marshaling puzzle: %v", err)
	}

	if err := os.MkdirAll("puzzles", 0755); err != nil {
		Log(LogLevelError, "Error creating puzzles directory: %v", err)
		return puzzle, fmt.Errorf("error creating puzzles directory: %v", err)
	}

	filename := fmt.Sprintf("puzzles/%s.json", puzzle.UUID)
	if err := ioutil.WriteFile(filename, data, 0644); err != nil {
		Log(LogLevelError, "Error writing puzzle file: %v", err)
		return puzzle, fmt.Errorf("error writing puzzle file: %v", err)
	}

	Log(LogLevelInfo, "Saved puzzle with UUID: %s", puzzle.UUID)
	return puzzle, nil
}

// ListPuzzles returns a list of saved puzzles
func ListPuzzles() ([]map[string]interface{}, error) {
	files, err := ioutil.ReadDir("puzzles")
	if err != nil {
		Log(LogLevelError, "Error reading puzzles directory: %v", err)
		return nil, err
	}

	var puzzles []map[string]interface{}
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") || strings.HasPrefix(file.Name(), ".DS_Store") {
			continue
		}

		// Get UUID without file extension
		uuid := file.Name()[:len(file.Name())-5]
		// Ensure uuid is long enough to slice; skip if invalid
		if len(uuid) < 8 {
			Log(LogLevelWarn, "Skipping invalid puzzle file with short UUID: %s", file.Name())
			continue
		}

		// Get create date from file
		date := file.ModTime().Format(time.RFC3339)

		// Load puzzle to get difficulty
		puzzle, err := LoadPuzzle(uuid)
		difficulty := 1 // Default difficulty
		if err == nil {
			// Use the stored difficulty value instead of recalculating
			difficulty = puzzle.Difficulty

			// If difficulty is not set (0 or invalid), calculate it based on filled cells
			if difficulty < 1 || difficulty > 9 {
				// Count the number of filled cells to estimate difficulty
				filledCells := 0
				for _, cell := range puzzle.Cells {
					if cell.Status == "s" {
						filledCells++
					}
				}

				// More filled cells = easier puzzle (lower difficulty)
				// Use floating point division for more accurate difficulty calculation
				difficulty = 10 - int(float64(filledCells)/9.0)
				if difficulty < 1 {
					difficulty = 1
				} else if difficulty > 9 {
					difficulty = 9
				}
			}
		}

		puzzles = append(puzzles, map[string]interface{}{
			"uuid":       uuid,
			"date":       date,
			"difficulty": difficulty,
			"shortId":    uuid[:8], // Safe to slice now
		})
	}

	// Sort puzzles by date (newest first)
	sort.Slice(puzzles, func(i, j int) bool {
		dateI, _ := time.Parse(time.RFC3339, puzzles[i]["date"].(string))
		dateJ, _ := time.Parse(time.RFC3339, puzzles[j]["date"].(string))
		return dateJ.Before(dateI) // Reverse order for newest first
	})

	Log(LogLevelDebug, "Listed %d puzzles", len(puzzles))
	return puzzles, nil
}

// LoadPuzzle loads a puzzle from disk by UUID
func LoadPuzzle(uuid string) (models.Puzzle, error) {
	var puzzle models.Puzzle

	filename := fmt.Sprintf("puzzles/%s.json", uuid)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		Log(LogLevelError, "Error reading puzzle file %s: %v", filename, err)
		return puzzle, err
	}

	if err := json.Unmarshal(data, &puzzle); err != nil {
		Log(LogLevelError, "Error unmarshaling puzzle: %v", err)
		return puzzle, err
	}

	Log(LogLevelDebug, "Loaded puzzle with UUID: %s", uuid)
	return puzzle, nil
}

// DeletePuzzle deletes a puzzle from disk by UUID
func DeletePuzzle(uuid string) error {
	filename := fmt.Sprintf("puzzles/%s.json", uuid)

	// Check if the file exists first
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		Log(LogLevelError, "Puzzle file %s does not exist", filename)
		return fmt.Errorf("puzzle with UUID %s does not exist", uuid)
	}

	// Delete the file
	if err := os.Remove(filename); err != nil {
		Log(LogLevelError, "Error deleting puzzle file %s: %v", filename, err)
		return fmt.Errorf("error deleting puzzle file: %v", err)
	}

	Log(LogLevelInfo, "Deleted puzzle with UUID: %s", uuid)
	return nil
}
