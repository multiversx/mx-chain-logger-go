package pipes

import (
	"os"

	logger "github.com/ElrondNetwork/elrond-go-logger"
)

// ChildMessenger is the messenger on child's part of the pipe
type ChildMessenger struct {
	Messenger
}

// NewChildMessenger creates a new messenger
func NewChildMessenger(profileReader *os.File, logsWriter *os.File) *ChildMessenger {
	receiver := NewReceiver(profileReader, &jsonMarshalizer{})
	sender := NewSender(logsWriter, &noopMarshalizer{})

	return &ChildMessenger{
		Messenger: *NewMessenger(receiver, sender),
	}
}

// ReceiveProfile reads an incoming profile
func (messenger *ChildMessenger) ReceiveProfile() (*logger.Profile, error) {
	profile := &logger.Profile{}
	err := messenger.Receive(profile)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

// SendLogLine sends a log line
func (messenger *ChildMessenger) SendLogLine(logLineMarshalized []byte) (int, error) {
	return messenger.Send(logLineMarshalized)
}
