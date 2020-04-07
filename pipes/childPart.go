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

	return &childPart{
		messenger:     messenger,
		outputSubject: outputSubject,
	}
}

func (part *childPart) StartLoop() error {
	err := part.addAsObserver()
	if err != nil {
		return err
	}

	go part.continuouslyReadProfile()
	return nil
}

func (part *childPart) addAsObserver() error {
	logLineFormatter, err := logger.NewLogLineWrapperFormatter(&jsonMarshalizer{})
	if err != nil {
		return err
	}

	part.outputSubject.ClearObservers()
	part.outputSubject.AddObserver(part, logLineFormatter)
	return nil
}

func (part *childPart) continuouslyReadProfile() {
	for {
		profile, err := part.messenger.ReceiveProfile()
		if err != nil {
			break
		}

		err = profile.Apply()
	}
}

func (part *childPart) Write(logLineMarshalized []byte) (int, error) {
	return part.messenger.SendLogLine(logLineMarshalized)
}
