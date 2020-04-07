package pipes

import (
	"os"

	logger "github.com/ElrondNetwork/elrond-go-logger"
)

const loggerSinkName = "loggerSink"

type parentPart struct {
	messenger          *ParentMessenger
	loggerSink         logger.Logger
	logLineMarshalizer logger.Marshalizer
	logsReader         *os.File
	logsWriter         *os.File
	profileReader      *os.File
	profileWriter      *os.File
}

// NewParentPart -
func NewParentPart(logLineMarshalizer logger.Marshalizer) (*parentPart, error) {
	loggerSink := logger.GetOrCreate(loggerSinkName)
	part := &parentPart{
		loggerSink:         loggerSink,
		logLineMarshalizer: logLineMarshalizer,
	}

	err := part.resetMessenger()
	if err != nil {
		return nil, err
	}

	return part, nil
}

func (part *parentPart) resetMessenger() error {
	err := part.resetPipes()
	if err != nil {
		return err
	}

	part.messenger = NewParentMessenger(part.logsReader, part.profileWriter, part.logLineMarshalizer)
	return nil
}

func (part *parentPart) resetPipes() error {
	closeFile(part.logsReader)
	closeFile(part.logsWriter)
	closeFile(part.profileReader)
	closeFile(part.profileWriter)

	var err error

	part.logsReader, part.logsWriter, err = os.Pipe()
	if err != nil {
		return err
	}

	part.profileReader, part.profileWriter, err = os.Pipe()
	if err != nil {
		return err
	}

	return nil
}

func closeFile(file *os.File) {
	if file != nil {
		_ = file.Close()
	}
}

func (part *parentPart) GetChildPipes() (*os.File, *os.File) {
	return part.profileReader, part.logsWriter
}

func (part *parentPart) StartLoop() {
	logger.SubscribeToProfileChange(part)
	part.forwardProfile()
	go part.continuouslyReadLogLines()
}

func (part *parentPart) OnProfileChanged() {
	part.forwardProfile()
}

func (part *parentPart) forwardProfile() {
	profile := logger.GetCurrentProfile()
	part.messenger.SendProfile(profile)
}

func (part *parentPart) continuouslyReadLogLines() {
	for {
		logLine, err := part.messenger.ReadLogLine()
		if err != nil {
			part.loggerSink.Error("continuouslyReadLogLines error", "err", err)
			break
		}

		part.loggerSink.Log(logLine)
	}
}

func (part *parentPart) Close() {
	logger.UnsubscribeFromProfileChange(part)
	// TODO: Also break loop, close pipes
}
