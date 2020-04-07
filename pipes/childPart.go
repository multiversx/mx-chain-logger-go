package pipes

import (
	"io"
	"os"

	logger "github.com/ElrondNetwork/elrond-go-logger"
)

var _ io.Writer = (*childPart)(nil)

type childPart struct {
	messenger     *ChildMessenger
	outputSubject logger.LogOutputHandler
}

// NewChildPart -
func NewChildPart(profileReader *os.File, logsWriter *os.File) *childPart {
	outputSubject := logger.GetLogOutputSubject()
	childPart := NewChildPartWithLogOutputSubject(outputSubject, profileReader, logsWriter)
	return childPart
}

// NewChildPartWithLogOutputSubject -
func NewChildPartWithLogOutputSubject(outputSubject logger.LogOutputHandler, profileReader *os.File, logsWriter *os.File) *childPart {
	messenger := NewChildMessenger(profileReader, logsWriter)

	part := &childPart{
		messenger:     messenger,
		outputSubject: outputSubject,
	}

	logLineFormatter, _ := logger.NewLogLineWrapperFormatter(&jsonMarshalizer{})
	outputSubject.ClearObservers()
	outputSubject.AddObserver(part, logLineFormatter)

	return part
}

func (part *childPart) Write(logLineMarshalized []byte) (int, error) {
	return part.messenger.Send(logLineMarshalized)
}
