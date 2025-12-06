package debug

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

// LogLevel corresponds to C++ LogLevel enum in voxelcore/src/debug/Logger.hpp
type LogLevel int

const (
	LogLevelPrint   LogLevel = iota // Print-debugging tool (printed without header)
	LogLevelDebug                   // Detailed information for debugging
	LogLevelInfo                    // General information about program progress
	LogLevelWarning                 // Potential issues that are not errors
	LogLevelError                   // Errors that prevent normal operation
)

// String returns the string representation of the LogLevel.
func (l LogLevel) String() string {
	switch l {
	case LogLevelPrint:
		return "PRINT"
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarning:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Global logger instance for file output
var (
	fileLogger     *log.Logger
	fileWriter     io.WriteCloser
	loggerInitOnce sync.Once
	loggerMutex    sync.Mutex // Protects fileLogger and fileWriter
)

// Init initializes the global logger to write to a file.
// Corresponds to C++ Logger::init(const std::string& filename)
func Init(filename string) {
	loggerInitOnce.Do(func() {
		// Open log file
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Failed to open log file %s: %v", filename, err)
		}
		fileWriter = file
		fileLogger = log.New(fileWriter, "", log.LstdFlags)
		log.SetOutput(io.MultiWriter(os.Stderr, fileWriter)) // Send to console and file
		log.Println("Logger initialized to file:", filename)
	})
}

// Flush writes any buffered log messages to the output.
// Corresponds to C++ Logger::flush()
func Flush() {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()
	if fileWriter != nil {
		if f, ok := fileWriter.(*os.File); ok {
			f.Sync() // Sync file to disk
		}
	}
}

// Logger corresponds to C++ debug::Logger class.
type Logger struct {
	name string
}

// NewLogger creates a new Logger instance with a given name.
// Corresponds to C++ Logger(const std::string& name) constructor.
func NewLogger(name string) *Logger {
	return &Logger{name: name}
}

// Log writes a message with a specific log level.
// Corresponds to C++ Logger::log(LogLevel level, std::string message)
func (l *Logger) Log(level LogLevel, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	
	// Prepend logger name and level for non-print messages
	if level != LogLevelPrint {
		message = fmt.Sprintf("%s [%s] %s", l.name, level.String(), message)
	}

	loggerMutex.Lock()
	defer loggerMutex.Unlock()

	// Always output to stderr for immediate visibility
	fmt.Fprintln(os.Stderr, message)

	// If file logger is initialized, write to file as well
	if fileLogger != nil {
		fileLogger.Output(2, message) // 2 means skip 2 call frames (Log, then actual call site)
	}
}

// Debug logs a debug message.
// Corresponds to C++ Logger::debug()
func (l *Logger) Debug(format string, args ...interface{}) {
	l.Log(LogLevelDebug, format, args...)
}

// Info logs an info message.
// Corresponds to C++ Logger::info()
func (l *Logger) Info(format string, args ...interface{}) {
	l.Log(LogLevelInfo, format, args...)
}

// Error logs an error message.
// Corresponds to C++ Logger::error()
func (l *Logger) Error(format string, args ...interface{}) {
	l.Log(LogLevelError, format, args...)
}

// Warning logs a warning message.
// Corresponds to C++ Logger::warning()
func (l *Logger) Warning(format string, args ...interface{}) {
	l.Log(LogLevelWarning, format, args...)
}

// Print logs a raw message without extra headers.
// Corresponds to C++ Logger::print()
func (l *Logger) Print(format string, args ...interface{}) {
	l.Log(LogLevelPrint, format, args...)
}