package pipes

import (
	"os"

	logger "github.com/ElrondNetwork/elrond-go-logger"
)

// ParentMessenger is the messenger on parent's part of the pipe
type ParentMessenger struct {
	Messenger
}

// NewParentMessenger creates a new messenger
func NewParentMessenger(logsReader *os.File, profileWriter *os.File) *ParentMessenger {
	receiver := NewReceiver(logsReader, &jsonMarshalizer{})
	sender := NewSender(profileWriter, &jsonMarshalizer{})

	return &ParentMessenger{
		Messenger: *NewMessenger(receiver, sender),
	}
}

// ReceiveLogLine reads a log line
func (messenger *ChildMessenger) ReceiveLogLine() (logger.LogLine, error) {
	logLine := logger.LogLine{}
	err := messenger.Receive(logLine)
	if err != nil {
		return logger.LogLine{}, err
	}

	return logLine, nil
}

// SendProfile sends a profile
func (messenger *ChildMessenger) SendProfile(profile logger.Profile) error {
	_, err := messenger.Send(profile)
	if err != nil {
		return err
	}

	return nil
}
