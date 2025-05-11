# Sudoku DJ

A Go-based Sudoku puzzle generator and solver with a RESTful API interface and frontend client.

## Quick Start

### Running the Application

1. **Start the server:**
   ```bash
   go build -o sudoku_dj ./cmd/app
   ./sudoku_dj --port 8081 --log-level info
   ```

2. **Start the client:**
   ```bash
   npm install
   npm start
   ```

3. Open your browser and navigate to `http://localhost:3000`

## Project Structure

```
sudoku_dj/
├── cmd/
│   └── app/
│       └── main.go        # Application entry point
├── internal/
│   ├── api/               # API handlers
│   │   └── handlers.go    # HTTP request handlers
│   ├── models/            # Data models
│   │   └── puzzle.go      # Sudoku puzzle model definition
│   ├── sudoku/            # Core sudoku logic
│   │   └── solver.go      # Puzzle generation and solving logic
│   └── utils/             # Utilities
│       ├── logging.go     # Logging system implementation
│       ├── storage.go     # File storage and persistence
│       └── utils.go       # General utility functions
├── c/                     # C implementation of Sudoku solver
│   ├── sudoku.c           # C implementation of solver
│   ├── sudoku_console.c   # Console interface for C solver
│   └── sudoku.h           # Header file with solver interface
├── puzzles/               # Directory where puzzles are stored
├── public/                # Static web assets
│   └── index.html         # Main HTML file for the web client
├── src/                   # Frontend source code
│   ├── App.js             # Main React component
│   ├── index.js           # React entry point
│   └── index.css          # Application styles
├── go.mod                 # Go module file
├── go.sum                 # Go dependency checksums
├── Dockerfile.server      # Dockerfile for the backend server
├── Dockerfile.client      # Dockerfile for the frontend client
├── docker-compose.yml     # Docker Compose configuration
├── package.json           # Node.js package configuration
└── package-lock.json      # Node.js dependency lock file
```

## Features

- Sudoku puzzle generation with configurable difficulty levels
- RESTful API for puzzle management
- Web-based UI for playing Sudoku
- Command-line interface for generating puzzles
- Grid display for visualizing puzzles in the terminal
- Flexible logging system with configurable log levels

## Getting Started

### Prerequisites

- Go 1.19 or higher
- GCC or compatible C compiler (for CGO)
- Node.js and npm (for frontend development)

### Building

```bash
# Build the backend application
go build -o sudoku_dj ./cmd/app

# Build the frontend
npm install
npm run build
```

### Running

```bash
# Run the application in server mode
./sudoku_dj --port 8081 --log-level info

# Generate a single puzzle without starting the server
./sudoku_dj --generate --difficulty 4
```

Command line options:
- `--port`: Port to run the server on (default: 8081)
- `--log-level`: Log level (error, warn, info, debug, trace)
- `--log-to-file`: Whether to log to a file in addition to stdout
- `--generate`: Generate a puzzle and exit without starting the server
- `--difficulty`: Difficulty level for puzzle generation (1-9, default: 5)

## API Endpoints

- `GET /` - Root endpoint, returns a simple HTML message
- `GET /sudoku` - Lists all available puzzles
- `POST /sudoku` - Generates a new Sudoku puzzle
  - Query parameters:
    - `difficulty` (1-9): Controls puzzle difficulty (default: 5)
    - `logLevel`: Controls logging level (default: "info")
- `GET /sudoku/{uuid}` - Retrieves a specific puzzle by UUID
- `POST /sudoku/validate` - Validates a puzzle solution
- `GET /sudoku/open?uuid={uuid}` - Opens a specific puzzle by UUID
- `POST /sudoku/save` - Saves a puzzle

## Puzzle Format

Puzzles are represented as JSON with the following structure:

```json
{
  "uuid": "unique-identifier",
  "createdAt": "ISO-8601-timestamp",
  "cells": {
    "01": { "value": 5, "notes": [], "status": "s" },
    "02": { "value": 0, "notes": [1, 2], "status": "" },
    ...
  }
}
```

Cell status values:
- `s`: System-generated (initial puzzle value)
- `u`: User-entered
- `c`: Correct (validated)
- `w`: Wrong (validated)

## Docker Support

The application includes Docker support for both backend and frontend:

```bash
# Build and run with Docker Compose
docker-compose up --build
```

**Disclaimer:** Docker configuration files may be out of date compared to the latest application code. For the most reliable experience, follow the manual setup instructions above.

## License

See the [LICENSE](LICENSE) file for details.
