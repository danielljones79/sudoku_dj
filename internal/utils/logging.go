package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// Log levels
const (
	LogLevelError = 1
	LogLevelWarn  = 2
	LogLevelInfo  = 3
	LogLevelDebug = 4
	LogLevelTrace = 5
)

// LogLevelName maps log level integers to their string representations
var LogLevelName = map[int]string{
	LogLevelError: "ERROR",
	LogLevelWarn:  "WARN",
	LogLevelInfo:  "INFO",
	LogLevelDebug: "DEBUG",
	LogLevelTrace: "TRACE",
}

var (
	currentLogLevel = LogLevelInfo // Default log level
	logMutex        sync.RWMutex
	logFile         *os.File
	logger          *log.Logger
	logToFile       bool
	logDirectory    = "logs"
)

// InitLogger initializes the logging system
func InitLogger(toFile bool) error {
	// Use standard logging before we set up our custom logger
	log.Println("[INFO] Initializing logging system")

	logMutex.Lock()
	defer logMutex.Unlock()

	// Set up logging to file if requested
	logToFile = toFile
	if logToFile {
		// Create log directory if it doesn't exist
		if err := os.MkdirAll(logDirectory, 0755); err != nil {
			log.Printf("[ERROR] Error creating log directory: %v", err)
			return fmt.Errorf("error creating log directory: %v", err)
		}

		// Create log file with current timestamp
		timestamp := time.Now().Format("2006-01-02_15-04-05")
		logFilePath := filepath.Join(logDirectory, fmt.Sprintf("sudoku_%s.log", timestamp))
		var err error
		logFile, err = os.Create(logFilePath)
		if err != nil {
			log.Printf("[ERROR] Error creating log file: %v", err)
			return fmt.Errorf("error creating log file: %v", err)
		}

		// Set up multiwriter to log to both stdout and file
		multiWriter := io.MultiWriter(os.Stdout, logFile)
		logger = log.New(multiWriter, "", log.LstdFlags)

		log.Printf("[INFO] Logging initialized, writing to %s", logFilePath)
	} else {
		// Log to stdout only
		logger = log.New(os.Stdout, "", log.LstdFlags)
		log.Println("[INFO] Logging initialized, writing to stdout only")
	}

	return nil
}

// CloseLogger properly closes the logger resources
func CloseLogger() {
	logMutex.Lock()
	defer logMutex.Unlock()

	if logFile != nil {
		Log(LogLevelInfo, "Closing log file")
		logFile.Close()
	}
}

// SetLogLevel sets the current log level
func SetLogLevel(level int) {
	logMutex.Lock()
	defer logMutex.Unlock()

	if level < LogLevelError {
		level = LogLevelError
	} else if level > LogLevelTrace {
		level = LogLevelTrace
	}

	oldLevel := currentLogLevel
	currentLogLevel = level

	if logger != nil {
		logger.Printf("[INFO] Log level changed from %s to %s",
			LogLevelName[oldLevel], LogLevelName[level])
	}
}

// GetLogLevel returns the current log level
func GetLogLevel() int {
	logMutex.RLock()
	defer logMutex.RUnlock()
	return currentLogLevel
}

// GetLogLevelName returns the name of the current log level
func GetLogLevelName() string {
	return LogLevelName[GetLogLevel()]
}

// LogLevelFromString converts a string log level to its integer value
func LogLevelFromString(level string) int {
	switch level {
	case "error":
		return LogLevelError
	case "warn":
		return LogLevelWarn
	case "info":
		return LogLevelInfo
	case "debug":
		return LogLevelDebug
	case "trace":
		return LogLevelTrace
	default:
		return LogLevelInfo
	}
}

// Log logs a message if the current log level is greater than or equal to the specified level
func Log(level int, format string, v ...interface{}) {
	logMutex.RLock()
	shouldLog := level <= currentLogLevel
	hasLogger := logger != nil
	logMutex.RUnlock()

	if shouldLog {
		prefix := ""
		if level >= LogLevelError && level <= LogLevelTrace {
			prefix = fmt.Sprintf("[%s] ", LogLevelName[level])
		}

		message := fmt.Sprintf(format, v...)
		if hasLogger {
			logger.Println(prefix + message)
		} else {
			// Fallback to standard log if logger isn't initialized
			log.Println(prefix + message)
		}
	}
}

// ParseLogLevel parses the log level from a string
func ParseLogLevel(logLevelStr string) string {
	if logLevelStr == "" {
		return ""
	}
	// Update the global log level based on the string
	logLevel := LogLevelFromString(logLevelStr)
	SetLogLevel(logLevel)

	return logLevelStr
}

// LogMem logs the current memory usage (useful for performance monitoring)
func LogMem(context string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	Log(LogLevelDebug, "%s - Memory stats: Alloc = %v MiB, TotalAlloc = %v MiB, Sys = %v MiB, NumGC = %v",
		context,
		bToMb(m.Alloc),
		bToMb(m.TotalAlloc),
		bToMb(m.Sys),
		m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
