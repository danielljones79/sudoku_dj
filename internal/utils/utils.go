package utils

import (
	"strconv"
)

// ParseDifficulty parses the difficulty level from a string
func ParseDifficulty(difficultyStr string) int {
	if difficultyStr == "" {
		return 5
	}
	difficulty, err := strconv.Atoi(difficultyStr)
	if err != nil || difficulty < 0 || difficulty > 9 {
		return 5
	}
	return difficulty
}
