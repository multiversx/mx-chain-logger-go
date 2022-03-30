package logger

import "io"

var DefaultLogLevel = &defaultLogLevel

// Observers -
func (los *logOutputSubject) Observers() ([]io.Writer, []Formatter) {
	los.mutObservers.RLock()
	defer los.mutObservers.RUnlock()

	return los.writers, los.formatters
}

// LogLevel -
func (l *logger) LogLevel() LogLevel {
	return l.logLevel
}

// IsASCII -
func IsASCII(data string) bool {
	return isASCII(data)
}
