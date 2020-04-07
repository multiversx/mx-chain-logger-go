package pipes

import (
	"os"
	"time"

	logger "github.com/ElrondNetwork/elrond-go-logger"
)

const loggerSinkName = "loggerSink"

type parentPart struct {
	messenger     *ParentMessenger
	loggerSink    logger.Logger
	logsReader    *os.File
	logsWriter    *os.File
	profileReader *os.File
	profileWriter *os.File
}

// NewParentPart -
func NewParentPart() (*parentPart, error) {
	loggerSink := logger.GetOrCreate(loggerSinkName)
	part := &parentPart{
		loggerSink: loggerSink,
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

	part.messenger = NewParentMessenger(part.logsReader, part.profileWriter)
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
	go part.continuouslyReadLogLines()
}

func (part *parentPart) continuouslyReadLogLines() {
	for {
		logLine, err := part.readLogLine()
		if err != nil {
			part.loggerSink.Error("continuouslyReadLogLines error", "err", err)
			break
		}

		part.loggerSink.Log(logLine)
	}
}

func (part *parentPart) readLogLine() (*logger.LogLine, error) {
	wrapper := &logger.LogLineWrapper{}
	part.messenger.Receive(wrapper)

	logLine := part.recoverLogLine(wrapper)
	return logLine, nil
}

func (part *parentPart) recoverLogLine(wrapper *logger.LogLineWrapper) *logger.LogLine {
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
