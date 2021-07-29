package logger

import (
	"io"

	"github.com/ElrondNetwork/elrond-go-logger/proto"
)

// Logger defines the behavior of a data logger component
type Logger interface {
	Trace(message string, args ...interface{})
	Debug(message string, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message string, args ...interface{})
	LogIfError(err error, args ...interface{})
	Log(line *LogLine)
	SetLevel(logLevel LogLevel)
	GetLevel() LogLevel
	IsInterfaceNil() bool
}

// LogLineHandler defines the get methods for a log line struct used by the formatter interface
type LogLineHandler interface {
	GetLoggerName() string
	GetCorrelation() proto.LogCorrelationMessage
	GetMessage() string
	GetLogLevel() int32
	GetArgs() []string
	GetTimestamp() int64
	IsInterfaceNil() bool
}

// Formatter describes what a log formatter should be able to do
type Formatter interface {
	Output(line LogLineHandler) []byte
	IsInterfaceNil() bool
}

// LogOutputHandler defines the properties of a subject-observer component
// able to output log lines
type LogOutputHandler interface {
	Output(line *LogLine)
	AddObserver(w io.Writer, format Formatter) error
	RemoveObserver(w io.Writer) error
	ClearObservers()
	IsInterfaceNil() bool
}

// Marshalizer defines the 2 basic operations: serialize (marshal) and deserialize (unmarshal)
type Marshalizer interface {
	Marshal(obj interface{}) ([]byte, error)
	Unmarshal(obj interface{}, buff []byte) error
	IsInterfaceNil() bool
}

// ProfileChangeObserver defines the interface for observing profile changes
type ProfileChangeObserver interface {
	OnProfileChanged()
}

// EpochStartNotifier defines the interface for epoch start notification
type EpochStartNotifier interface {
	RegisterForEpochChangeConfirmed(handler func(epoch uint32))
	IsInterfaceNil() bool
}

// LogLifeSpanner defines a notification channel for the file logging lifespan
type LogLifeSpanner interface {
	GetNotification() <-chan string
	IsInterfaceNil() bool
	SetCurrentFile(string)
}

// LogLifeSpanFactory defines the methods for creating a log lifeSpanner
type LogLifeSpanFactory interface {
	CreateLogLifeSpanner(args LogLifeSpanFactoryArgs) (LogLifeSpanner, error)
}

// FileSizeCheckHandler defines the method needed for getting a file size
type FileSizeCheckHandler interface {
	GetSize(path string) (int64, error)
	IsInterfaceNil() bool
}

type LifeSpannerWrapper interface {
	SetLifeSpanner(spanner LogLifeSpanner)
	SetCurrentFile(newFile string)
	GetNotificationChannel() <-chan string
}
