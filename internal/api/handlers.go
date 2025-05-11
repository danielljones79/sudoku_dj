package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/danjones/sudoku_dj/internal/models"
	"github.com/danjones/sudoku_dj/internal/sudoku"
	"github.com/danjones/sudoku_dj/internal/utils"
)

// EnableCORS adds CORS headers to support cross-origin requests
func EnableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// HandleRoot serves the root endpoint
func HandleRoot(w http.ResponseWriter, r *http.Request) {
	utils.Log(utils.LogLevelInfo, "Handling request to root endpoint")

	EnableCORS(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "<h1>Sudoku Server</h1><p>Server is running. Visit <a href='/sudoku'>/sudoku</a> to get started.</p>")
}

// HandleSudokuRequest handles requests to the /sudoku endpoint
func HandleSudokuRequest(w http.ResponseWriter, r *http.Request) {
	utils.Log(utils.LogLevelInfo, "Handling request to /sudoku endpoint: %s %s", r.Method, r.URL.Path)

	EnableCORS(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Parse URL path to determine what to serve
	pathParts := strings.Split(r.URL.Path, "/")

	// If path is just /sudoku or /sudoku/
	if len(pathParts) <= 2 || pathParts[2] == "" {
		if r.Method == "POST" {
			// Generate new puzzle
			HandleGenerateSudokuPuzzle(w, r)
		} else {
			// List puzzles
			HandleListPuzzles(w, r)
		}
		return
	}

	// Handle requests to /sudoku/{uuid}
	HandlePuzzleByUUID(w, r, pathParts[2])
}

// HandleGenerateSudokuPuzzle generates a new Sudoku puzzle
func HandleGenerateSudokuPuzzle(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	err := r.ParseForm()
	if err != nil {
		utils.Log(utils.LogLevelError, "Failed to parse form data: %v", err)
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	difficulty := utils.ParseDifficulty(r.FormValue("difficulty"))
	logLevel := utils.ParseLogLevel(r.FormValue("log_level"))

	// Set log level if provided
	if logLevel != "" {
		oldLevel := utils.GetLogLevel()
		utils.SetLogLevel(utils.LogLevelFromString(logLevel))
		utils.Log(utils.LogLevelInfo, "Log level changed from %d to %d for this request", oldLevel, utils.GetLogLevel())
	}

	utils.Log(utils.LogLevelInfo, "Generating new puzzle with difficulty: %d", difficulty)

	// Generate puzzle
	puzzle := sudoku.CreatePuzzle(difficulty, logLevel)

	// Save puzzle to disk
	savedPuzzle, err := utils.SavePuzzle(puzzle)
	if err != nil {
		utils.Log(utils.LogLevelError, "Failed to save puzzle: %v", err)
		http.Error(w, "Failed to save puzzle", http.StatusInternalServerError)
		return
	}

	// Return puzzle
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(savedPuzzle)

	utils.Log(utils.LogLevelInfo, "Successfully generated and saved puzzle with UUID: %s", savedPuzzle.UUID)
}

// HandleListPuzzles lists all available puzzles
func HandleListPuzzles(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	err := r.ParseForm()
	if err != nil {
		utils.Log(utils.LogLevelError, "Failed to parse form data: %v", err)
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	logLevel := utils.ParseLogLevel(r.FormValue("log_level"))

	// Set log level if provided
	if logLevel != "" {
		oldLevel := utils.GetLogLevel()
		utils.SetLogLevel(utils.LogLevelFromString(logLevel))
		utils.Log(utils.LogLevelInfo, "Log level changed from %d to %d for this request", oldLevel, utils.GetLogLevel())
	}

	utils.Log(utils.LogLevelInfo, "Listing available puzzles")

	// Get puzzles
	puzzles, err := utils.ListPuzzles()
	if err != nil {
		utils.Log(utils.LogLevelError, "Failed to list puzzles: %v", err)
		http.Error(w, "Failed to list puzzles", http.StatusInternalServerError)
		return
	}

	// Return puzzles
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(puzzles)

	utils.Log(utils.LogLevelInfo, "Successfully listed %d puzzles", len(puzzles))
}

// HandleOpenPuzzle opens a specific puzzle by UUID
func HandleOpenPuzzle(w http.ResponseWriter, r *http.Request, uuid string) {
	// Parse query parameters
	err := r.ParseForm()
	if err != nil {
		utils.Log(utils.LogLevelError, "Failed to parse form data: %v", err)
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	logLevel := utils.ParseLogLevel(r.FormValue("log_level"))

	// Set log level if provided
	if logLevel != "" {
		oldLevel := utils.GetLogLevel()
		utils.SetLogLevel(utils.LogLevelFromString(logLevel))
		utils.Log(utils.LogLevelInfo, "Log level changed from %d to %d for this request", oldLevel, utils.GetLogLevel())
	}

	utils.Log(utils.LogLevelInfo, "Opening puzzle with UUID: %s", uuid)

	// Load puzzle
	puzzle, err := utils.LoadPuzzle(uuid)
	if err != nil {
		utils.Log(utils.LogLevelError, "Failed to load puzzle %s: %v", uuid, err)
		http.Error(w, "Failed to load puzzle", http.StatusNotFound)
		return
	}

	// Return puzzle
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(puzzle)

	utils.Log(utils.LogLevelInfo, "Successfully loaded puzzle with UUID: %s", uuid)
}

// HandleValidatePuzzle validates a puzzle solution
func HandleValidatePuzzle(w http.ResponseWriter, r *http.Request, uuid string) {
	// Parse query parameters
	err := r.ParseForm()
	if err != nil {
		utils.Log(utils.LogLevelError, "Failed to parse form data: %v", err)
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	logLevel := utils.ParseLogLevel(r.FormValue("log_level"))

	// Set log level if provided
	if logLevel != "" {
		oldLevel := utils.GetLogLevel()
		utils.SetLogLevel(utils.LogLevelFromString(logLevel))
		utils.Log(utils.LogLevelInfo, "Log level changed from %d to %d for this request", oldLevel, utils.GetLogLevel())
	}

	utils.Log(utils.LogLevelInfo, "Validating puzzle with UUID: %s", uuid)

	// Decode puzzle from request body
	var puzzle models.Puzzle
	err = json.NewDecoder(r.Body).Decode(&puzzle)
	if err != nil {
		utils.Log(utils.LogLevelError, "Failed to decode puzzle from request body: %v", err)
		http.Error(w, "Failed to decode puzzle from request body", http.StatusBadRequest)
		return
	}

	// Validate puzzle solution
	validatedPuzzle, _ := sudoku.ValidateSolution(puzzle)

	// Return validated puzzle
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(validatedPuzzle)

	utils.Log(utils.LogLevelInfo, "Successfully validated puzzle with UUID: %s", uuid)
}

// HandleSavePuzzle saves a puzzle to the server
func HandleSavePuzzle(w http.ResponseWriter, r *http.Request, uuid string) {
	// Parse query parameters
	err := r.ParseForm()
	if err != nil {
		utils.Log(utils.LogLevelError, "Failed to parse form data: %v", err)
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	logLevel := utils.ParseLogLevel(r.FormValue("log_level"))

	// Set log level if provided
	if logLevel != "" {
		oldLevel := utils.GetLogLevel()
		utils.SetLogLevel(utils.LogLevelFromString(logLevel))
		utils.Log(utils.LogLevelInfo, "Log level changed from %d to %d for this request", oldLevel, utils.GetLogLevel())
	}

	utils.Log(utils.LogLevelInfo, "Saving puzzle with UUID: %s", uuid)

	// Decode puzzle from request body
	var puzzle models.Puzzle
	err = json.NewDecoder(r.Body).Decode(&puzzle)
	if err != nil {
		utils.Log(utils.LogLevelError, "Failed to decode puzzle from request body: %v", err)
		http.Error(w, "Failed to decode puzzle from request body", http.StatusBadRequest)
		return
	}

	// If new save, update creation time
	existingPuzzle, err := utils.LoadPuzzle(uuid)
	if err != nil {
		puzzle.CreatedAt = time.Now().Format(time.RFC3339)
	} else {
		puzzle.CreatedAt = existingPuzzle.CreatedAt
	}

	// Save puzzle
	savedPuzzle, err := utils.SavePuzzle(puzzle)
	if err != nil {
		utils.Log(utils.LogLevelError, "Failed to save puzzle: %v", err)
		http.Error(w, "Failed to save puzzle", http.StatusInternalServerError)
		return
	}

	// Return saved puzzle
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(savedPuzzle)

	utils.Log(utils.LogLevelInfo, "Successfully saved puzzle with UUID: %s", uuid)
}

// HandlePuzzleByUUID handles requests to /sudoku/{uuid}
func HandlePuzzleByUUID(w http.ResponseWriter, r *http.Request, uuid string) {
	utils.Log(utils.LogLevelDebug, "Handling request to /sudoku/%s: %s", uuid, r.Method)

	switch r.Method {
	case "GET":
		HandleOpenPuzzle(w, r, uuid)
	case "POST":
		HandleValidatePuzzle(w, r, uuid)
	case "PUT":
		HandleSavePuzzle(w, r, uuid)
	case "DELETE":
		HandleDeletePuzzle(w, r, uuid)
	default:
		utils.Log(utils.LogLevelWarn, "Unsupported method %s for /sudoku/%s", r.Method, uuid)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleDeletePuzzle deletes a puzzle from the server
func HandleDeletePuzzle(w http.ResponseWriter, r *http.Request, uuid string) {
	// Parse query parameters
	err := r.ParseForm()
	if err != nil {
		utils.Log(utils.LogLevelError, "Failed to parse form data: %v", err)
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	logLevel := utils.ParseLogLevel(r.FormValue("log_level"))

	// Set log level if provided
	if logLevel != "" {
		oldLevel := utils.GetLogLevel()
		utils.SetLogLevel(utils.LogLevelFromString(logLevel))
		utils.Log(utils.LogLevelInfo, "Log level changed from %d to %d for this request", oldLevel, utils.GetLogLevel())
	}

	utils.Log(utils.LogLevelInfo, "Deleting puzzle with UUID: %s", uuid)

	// Delete puzzle
	err = utils.DeletePuzzle(uuid)
	if err != nil {
		utils.Log(utils.LogLevelError, "Failed to delete puzzle: %v", err)
		if strings.Contains(err.Error(), "does not exist") {
			http.Error(w, "Puzzle not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to delete puzzle", http.StatusInternalServerError)
		}
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"status":  "success",
		"message": fmt.Sprintf("Puzzle %s deleted successfully", uuid),
	}
	json.NewEncoder(w).Encode(response)

	utils.Log(utils.LogLevelInfo, "Successfully deleted puzzle with UUID: %s", uuid)
}

// SetupRoutes configures all API routes
func SetupRoutes() http.Handler {
	utils.Log(utils.LogLevelInfo, "Setting up API routes")

	mux := http.NewServeMux()
	mux.HandleFunc("/", HandleRoot)
	mux.HandleFunc("/sudoku", HandleSudokuRequest)
	mux.HandleFunc("/sudoku/", HandleSudokuRequest)

	utils.Log(utils.LogLevelInfo, "API routes configured successfully")
	return mux
}
