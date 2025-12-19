package logger

import (
	"strings"
	"time"

	"github.com/abaxoth0/Vega/libs/go/packages/structs"
)

type logLevel uint8

const (
	TraceLogLevel logLevel = iota
	DebugLogLevel
	InfoLogLevel
	WarningLogLevel
	ErrorLogLevel
	// Logs with this level will be handled immediately after calling Log().
	// Logger will call os.Exit(1) after this log.
	FatalLogLevel
	// Logs with this level will be handled immediately after calling Log().
	// Logger will call os.Exit(1) after this log.
	PanicLogLevel
)

var logLevelToStrMap = map[logLevel]string{
	TraceLogLevel:   "TRACE",
	DebugLogLevel:   "DEBUG",
	InfoLogLevel:    "INFO",
	WarningLogLevel: "WARNING",
	ErrorLogLevel:   "ERROR",
	FatalLogLevel:   "FATAL",
	PanicLogLevel:   "PANIC",
}

func (s logLevel) String() string {
	return logLevelToStrMap[s]
}

type LogEntry struct {
	rawLevel  logLevel
	Timestamp time.Time 	`json:"ts"`
	Service   string    	`json:"service"`
	Instance  string    	`json:"instance"`
	Level     string 		`json:"level"`
	Source    string 		`json:"source,omitempty"`
	Message   string 		`json:"msg"`
	Error     string 		`json:"error,omitempty"`
	Meta      structs.Meta  `json:"meta,omitempty"`
}

var (
	serviceName 	string = "undefined"
	serviceInstance string = "undefined"
)

func SetServiceName(name string) {
	name = strings.Trim(name, " \r\n")
	if name == "" {
		return
	}
	serviceName = name
}

func GetServiceName() string {
	return serviceName
}

func SetServiceInstance(instance string) {
	instance = strings.Trim(instance, " \r\n")
	if instance == "" {
		return
	}
	serviceInstance = instance
}

func GetServiceInstance() string {
	return serviceInstance
}

// Creates a new log entry. Timestamp is time.Now().
// If level is not error, fatal or panic, then Error will be empty, even if err specified.
func NewLogEntry(level logLevel, src string, msg string, err string, meta structs.Meta) LogEntry {
	e := LogEntry{
		rawLevel:  level,
		Timestamp: time.Now(),
		Service:   serviceName,
		Instance:  serviceInstance,
		Level:     level.String(),
		Source:    src,
		Message:   msg,
		Meta:      meta,
	}

	// error, fatal, panic
	if level >= ErrorLogLevel {
		e.Error = err
	}

	return e
}
