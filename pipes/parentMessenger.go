package pipes

import (
	"encoding/json"
	"os"
	"time"

	logger "github.com/ElrondNetwork/elrond-go-logger"
)

// ParentMessenger is the messenger on parent's part of the pipe
type ParentMessenger struct {
	Messenger
}

// NewParentMessenger creates a new messenger
func NewParentMessenger(logsReader *os.File, profileWriter *os.File) *ParentMessenger {
	receiver := NewReceiver(logsReader)
	sender := NewSender(profileWriter)

	return &ParentMessenger{
		Messenger: *NewMessenger(receiver, sender),
	}
}

// ReceiveLogLine reads a log line
func (messenger *ParentMessenger) ReceiveLogLine() (*logger.LogLine, error) {
	buffer, err := messenger.Receive()
	if err != nil {
		return nil, err
	}

	wrapper := &logger.LogLineWrapper{}
	err = json.Unmarshal(buffer, wrapper)
	if err != nil {
		return nil, err
	}

	logLine := messenger.recoverLogLine(wrapper)
	return logLine, nil
}

func (messenger *ParentMessenger) recoverLogLine(wrapper *logger.LogLineWrapper) *logger.LogLine {
	logLine := &logger.LogLine{
		LoggerName:  wrapper.LoggerName,
		Correlation: wrapper.Correlation,
		Message:     wrapper.Message,
		LogLevel:    logger.LogLevel(wrapper.LogLevel),
		Args:        make([]interface{}, len(wrapper.Args)),
		Timestamp:   time.Unix(0, wrapper.Timestamp),
	}

	for i, str := range wrapper.Args {
		logLine.Args[i] = str
	}

	return logLine
}

// SendProfile sends a profile
func (messenger *ParentMessenger) SendProfile(profile logger.Profile) error {
	buffer, err := profile.Marshal()
	if err != nil {
		return err
	}

	_, err = messenger.Send(buffer)
	if err != nil {
		return err
	}

	return nil
}
