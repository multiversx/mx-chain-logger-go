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
	receiver := NewReceiver(profileReader)
	sender := NewSender(logsWriter)

	return &ChildMessenger{
		Messenger: *NewMessenger(receiver, sender),
	}
}

// ReceiveProfile reads an incoming profile
func (messenger *ChildMessenger) ReceiveProfile() (logger.Profile, error) {
	buffer, err := messenger.Receive()
	if err != nil {
		return logger.Profile{}, err
	}

	return logger.UnmarshalProfile(buffer)
}

// SendLogLine sends a log line
func (messenger *ChildMessenger) SendLogLine(logLineMarshalized []byte) (int, error) {
	return messenger.Send(logLineMarshalized)
}
