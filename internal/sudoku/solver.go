package sudoku

/*
#cgo CFLAGS: -g -Wall
#cgo LDFLAGS: -L../../c -lsudoku
#include <stdio.h>
#include <stdbool.h>
#include "../../sudoku.h"
*/
import "C"
import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
	"unsafe"

	"github.com/google/uuid"

	"github.com/danjones/sudoku_dj/internal/models"
	"github.com/danjones/sudoku_dj/internal/utils"
)

// PuzzleGrid represents a Sudoku grid
type PuzzleGrid [9][9]C.int

// InitSolver initializes the sudoku solver
func InitSolver() {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())
	utils.Log(utils.LogLevelInfo, "Sudoku solver initialized")
}

// PrintGrid returns a string representation of the grid for logging
func PrintGrid(grid PuzzleGrid) string {
	var output string

	// Top border
	output += "+-------+-------+-------+\n"

	for row := 0; row < 9; row++ {
		output += "| "
		for col := 0; col < 9; col++ {
			// Print number or space if zero
			val := grid[row][col]
			if val == 0 {
				output += ". "
			} else {
				output += fmt.Sprintf("%d ", val)
			}

			// Add vertical dividers
			if col == 2 || col == 5 {
				output += "| "
			}
		}
		output += "|\n"

		// Add horizontal dividers
		if row == 2 || row == 5 {
			output += "+-------+-------+-------+\n"
		}
	}

	// Bottom border
	output += "+-------+-------+-------+"
	return output
}

// CreatePuzzle generates a new Sudoku puzzle with the specified difficulty level
func CreatePuzzle(level int, logLevel string) models.Puzzle {
	utils.Log(utils.LogLevelInfo, "Creating new puzzle with difficulty level %d", level)

	var emptyGrid PuzzleGrid
	initializeGrid(&emptyGrid)
	placeRandomNumbers(&emptyGrid)

	// Log the initial grid with random numbers
	utils.Log(utils.LogLevelDebug, "Initial grid with random seeds:\n%s", PrintGrid(emptyGrid))

	// Convert to cells and attempt to solve
	cells := gridToCells(emptyGrid)
	utils.Log(utils.LogLevelDebug, "Attempting to solve initial grid")
	solvedCells, solved := AttemptSolve(cells)

	if !solved {
		utils.Log(utils.LogLevelError, "Failed to solve the initial grid")
		return models.Puzzle{Cells: make(map[string]models.Cell)}
	}

	// Convert back to grid for refinement
	solutionGrid := cellsToGrid(solvedCells)

	// Log the solved grid
	utils.Log(utils.LogLevelDebug, "Solved grid:\n%s", PrintGrid(solutionGrid))

	utils.Log(utils.LogLevelDebug, "Refining puzzle to difficulty level %d", level)
	refinedGrid := refinePuzzle(solutionGrid, level)

	// Log the final grid with cells removed
	utils.Log(utils.LogLevelDebug, "Final puzzle grid (with cells removed):\n%s", PrintGrid(refinedGrid))

	// Create puzzle with system cells marked
	puzzle := models.Puzzle{
		UUID:       uuid.New().String(),
		Cells:      gridToCells(refinedGrid),
		CreatedAt:  time.Now().Format(time.RFC3339),
		Difficulty: level,
	}

	// Mark system-generated cells
	nonEmptyCells := 0
	for posKey, cell := range puzzle.Cells {
		if cell.Value != 0 {
			cell.Status = "s"
			puzzle.Cells[posKey] = cell
			nonEmptyCells++
		}
	}

	utils.Log(utils.LogLevelInfo, "Created puzzle with %d filled cells", nonEmptyCells)
	return puzzle
}

// AttemptSolve attempts to solve a Sudoku puzzle
func AttemptSolve(cells map[string]models.Cell) (map[string]models.Cell, bool) {
	utils.Log(utils.LogLevelDebug, "Attempting to solve puzzle")
	startTime := time.Now()

	grid := cellsToGrid(cells)
	solved := C.solve_sudoku((*C.int)(unsafe.Pointer(&grid[0][0])))

	duration := time.Since(startTime)
	if bool(solved) {
		utils.Log(utils.LogLevelInfo, "Puzzle solved successfully in %v", duration)
		return gridToCells(grid), true
	}

	utils.Log(utils.LogLevelWarn, "Failed to solve puzzle after %v", duration)
	return cells, false
}

// ValidateSolution validates user-entered cells against a solution
func ValidateSolution(puzzle models.Puzzle) (models.Puzzle, bool) {
	utils.Log(utils.LogLevelDebug, "Validating puzzle solution")

	// Create a grid with only system cells
	var systemGrid PuzzleGrid
	systemCellCount := 0
	for posKey, cell := range puzzle.Cells {
		if cell.Status == "s" {
			pos, _ := strconv.Atoi(posKey)
			row, col := (pos-1)/9, (pos-1)%9
			systemGrid[row][col] = C.int(cell.Value)
			systemCellCount++
		}
	}

	utils.Log(utils.LogLevelDebug, "System cells grid:\n%s", PrintGrid(systemGrid))
	utils.Log(utils.LogLevelDebug, "Solving grid with %d system cells", systemCellCount)

	// Get the solution using only system cells
	startTime := time.Now()
	solved := C.solve_sudoku((*C.int)(unsafe.Pointer(&systemGrid[0][0])))
	duration := time.Since(startTime)

	if !bool(solved) {
		utils.Log(utils.LogLevelWarn, "Failed to solve with system cells after %v", duration)
		return puzzle, false
	}

	utils.Log(utils.LogLevelDebug, "Solved grid:\n%s", PrintGrid(systemGrid))
	utils.Log(utils.LogLevelDebug, "Solved with system cells in %v", duration)

	// Create a grid with user's solution for display
	var userGrid PuzzleGrid
	for posKey, cell := range puzzle.Cells {
		pos, _ := strconv.Atoi(posKey)
		row, col := (pos-1)/9, (pos-1)%9
		userGrid[row][col] = C.int(cell.Value)
	}
	utils.Log(utils.LogLevelDebug, "User's solution grid:\n%s", PrintGrid(userGrid))

	// Validate each user-entered cell against the solution
	correctCount := 0
	wrongCount := 0
	for posKey, cell := range puzzle.Cells {
		// Skip system cells and empty cells
		if cell.Status == "s" || cell.Value == 0 {
			continue
		}

		// Check if this cell's value matches the solution
		pos, _ := strconv.Atoi(posKey)
		row, col := (pos-1)/9, (pos-1)%9
		if systemGrid[row][col] == C.int(cell.Value) {
			cell.Status = "c" // Correct
			correctCount++
		} else {
			cell.Status = "w" // Wrong
			wrongCount++
		}
		puzzle.Cells[posKey] = cell
	}

	utils.Log(utils.LogLevelInfo, "Validation complete: %d correct, %d wrong cells", correctCount, wrongCount)
	return puzzle, true
}

// initializeGrid initializes a grid with zeros
func initializeGrid(grid *PuzzleGrid) {
	utils.Log(utils.LogLevelTrace, "Initializing empty grid")
	for row := range grid {
		for col := range grid[row] {
			grid[row][col] = 0
		}
	}
}

// placeRandomNumbers places random numbers in the grid
func placeRandomNumbers(grid *PuzzleGrid) {
	utils.Log(utils.LogLevelDebug, "Placing initial random numbers")
	for num := 1; num <= 9; num++ {
		x, y := rand.Intn(9), rand.Intn(9)
		for grid[x][y] != 0 {
			x, y = rand.Intn(9), rand.Intn(9)
		}
		grid[x][y] = C.int(num)
		utils.Log(utils.LogLevelTrace, "Placed %d at position (%d,%d)", num, x, y)
	}
}

// refinePuzzle refines a solved puzzle to create a playable puzzle
func refinePuzzle(grid PuzzleGrid, difficulty int) PuzzleGrid {
	utils.Log(utils.LogLevelDebug, "Refining puzzle to achieve difficulty level %d", difficulty)

	var refinedGrid = grid
	attempts := 0
	removedCount := 0

	// Create a list of all positions to check in random order
	positions := make([][2]int, 81)
	for i := 0; i < 81; i++ {
		positions[i] = [2]int{i / 9, i % 9}
	}

	// This creates a random permutation of all cell positions
	for i := 0; i < len(positions); i++ {
		j := i + rand.Intn(len(positions)-i)
		positions[i], positions[j] = positions[j], positions[i]
	}

	// Try to remove each cell once
	for _, pos := range positions {
		x, y := pos[0], pos[1]
		attempts++

		// Skip already empty cells
		if refinedGrid[x][y] == 0 {
			continue
		}

		removedValue := refinedGrid[x][y]
		refinedGrid[x][y] = 0

		// Check if removal maintains a unique solution
		if newCnt := int(C.count_solutions((*C.int)(unsafe.Pointer(&refinedGrid[0][0])))); newCnt > 1 {
			// Restore value if it creates multiple solutions
			refinedGrid[x][y] = removedValue
			utils.Log(utils.LogLevelTrace, "Removing cell (%d,%d) would create multiple solutions, restoring", x, y)
		} else {
			// Successful removal
			removedCount++
			utils.Log(utils.LogLevelTrace, "Removed value %d at position (%d,%d)", removedValue, x, y)
		}

		// Log every 10th attempt to show progress with current grid state
		if attempts%10 == 0 {
			utils.Log(utils.LogLevelTrace, "Current grid after %d removal attempts, %d removals:\n%s",
				attempts, removedCount, PrintGrid(refinedGrid))
		}
	}

	utils.Log(utils.LogLevelDebug, "Initial refinement complete after %d attempts with %d cells removed", attempts, removedCount)

	// Add some values back based on difficulty
	addedBack := 0
	for i := 9 - difficulty; i >= 0; i-- {
		x, y := rand.Intn(9), rand.Intn(9)
		for refinedGrid[x][y] != 0 {
			x, y = rand.Intn(9), rand.Intn(9)
		}
		refinedGrid[x][y] = grid[x][y]
		addedBack++
		utils.Log(utils.LogLevelTrace, "Added back value %d at position (%d,%d)", grid[x][y], x, y)
	}

	utils.Log(utils.LogLevelDebug, "Added back %d values based on difficulty level", addedBack)

	return refinedGrid
}

// gridToCells converts a PuzzleGrid to a map of cells
func gridToCells(grid PuzzleGrid) map[string]models.Cell {
	utils.Log(utils.LogLevelTrace, "Converting grid to cells")
	cells := make(map[string]models.Cell)
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			pos := row*9 + col + 1
			posKey := fmt.Sprintf("%02d", pos)
			cells[posKey] = models.Cell{
				Value:  int(grid[row][col]),
				Notes:  []int{},
				Status: "", // Will be set to "s" for system-generated cells
			}
		}
	}
	return cells
}

// cellsToGrid converts a map of cells to a PuzzleGrid
func cellsToGrid(cells map[string]models.Cell) PuzzleGrid {
	utils.Log(utils.LogLevelTrace, "Converting cells to grid")
	var grid PuzzleGrid
	for posKey, cell := range cells {
		pos, _ := strconv.Atoi(posKey)
		row, col := (pos-1)/9, (pos-1)%9
		grid[row][col] = C.int(cell.Value)
	}
	return grid
}
