package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/danjones/sudoku_dj/internal/api"
	"github.com/danjones/sudoku_dj/internal/sudoku"
	"github.com/danjones/sudoku_dj/internal/utils"
)

func main() {
	// Print startup message
	log.Println("Starting Sudoku DJ application...")

	// Parse command line flags
	port := flag.String("port", "8081", "Port to run the server on")
	logLevel := flag.String("log-level", "info", "Log level (error, warn, info, debug, trace)")
	logToFile := flag.Bool("log-to-file", false, "Whether to log to a file in addition to stdout")

	// Show usage if help flag is present
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	// Echo back the parameters we're using
	log.Printf("Command-line parameters:")
	log.Printf("  - port: %s", *port)
	log.Printf("  - log-level: %s", *logLevel)
	log.Printf("  - log-to-file: %v", *logToFile)

	// Initialize the logging system
	log.Println("Initializing logging system...")
	if err := utils.InitLogger(*logToFile); err != nil {
		log.Fatalf("Error initializing logger: %v", err)
	}
	// Make sure to close log file when application exits
	defer utils.CloseLogger()

	// Set the global log level from command line flag
	if *logLevel != "" {
		levelInt := utils.LogLevelFromString(strings.ToLower(*logLevel))
		utils.SetLogLevel(levelInt)
		utils.Log(utils.LogLevelInfo, "Set log level to %s (%d)", *logLevel, levelInt)
	}

	utils.Log(utils.LogLevelInfo, "Starting Sudoku DJ with log level: %s", utils.GetLogLevelName())

	// Create puzzles directory if it doesn't exist
	utils.Log(utils.LogLevelInfo, "Checking if puzzles directory exists...")
	if err := os.MkdirAll("puzzles", 0755); err != nil {
		utils.Log(utils.LogLevelError, "Error creating puzzles directory: %v", err)
		os.Exit(1)
	}

	// Initialize sudoku solver
	utils.Log(utils.LogLevelInfo, "Initializing Sudoku solver...")
	sudoku.InitSolver()

	// Setup routes
	utils.Log(utils.LogLevelInfo, "Setting up API routes...")
	handler := api.SetupRoutes()

	// Start server
	utils.Log(utils.LogLevelInfo, "Starting Sudoku DJ on port %s...", *port)
	utils.Log(utils.LogLevelInfo, "Server is ready to accept connections")
	utils.Log(utils.LogLevelInfo, "Press Ctrl+C to stop the server")

	if err := http.ListenAndServe(":"+*port, handler); err != nil {
		utils.Log(utils.LogLevelError, "Server error: %v", err)
		os.Exit(1)
	}
}
