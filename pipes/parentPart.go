package pipes

import (
	"bufio"
	"io"
	"os"
	"strings"

	logger "github.com/ElrondNetwork/elrond-go-logger"
)

const logLinesSinkName = "logLinesSink"
const textOutputSinkName = "textOutputSink"

type parentPart struct {
	messenger          *ParentMessenger
	logLinesSink       logger.Logger
	textOutputSink     logger.Logger
	logLineMarshalizer logger.Marshalizer
	logsReader         *os.File
	logsWriter         *os.File
	profileReader      *os.File
	profileWriter      *os.File
}

// NewParentPart creates a new logs receiver part (in the parent process)
func NewParentPart(logLineMarshalizer logger.Marshalizer) (*parentPart, error) {
	part := &parentPart{
		logLinesSink:       logger.GetOrCreate(logLinesSinkName),
		textOutputSink:     logger.GetOrCreate(textOutputSinkName),
		logLineMarshalizer: logLineMarshalizer,
	}

	err := part.initializePipes()
	if err != nil {
		return nil, err
	}

	part.initializeMessenger()

	return part, nil
}

func (part *parentPart) initializePipes() error {
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

func (part *parentPart) initializeMessenger() {
	part.messenger = NewParentMessenger(part.logsReader, part.profileWriter, part.logLineMarshalizer)
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
			part.logLinesSink.Error("continuouslyReadLogLines error", "err", err)
			break
		}

		part.logLinesSink.Log(logLine)
	}
}

func (part *parentPart) ContinuouslyReadTextualOutput(childStdout io.Reader, childStderr io.Reader, tag string) {
	stdoutReader := bufio.NewReader(childStdout)
	stderrReader := bufio.NewReader(childStderr)

	go func() {
		for {
			line, err := stdoutReader.ReadString('\n')
			if err != nil {
				break
			}

			line = strings.TrimSpace(line)
			part.textOutputSink.Trace(tag, "line", line)
		}
	}()

	go func() {
		for {
			line, err := stderrReader.ReadString('\n')
			if err != nil {
				break
			}

			line = strings.TrimSpace(line)
			part.textOutputSink.Error(tag, "line", line)
		}
	}()
}

func (part *parentPart) Close() {
	logger.UnsubscribeFromProfileChange(part)

	_ = part.logsReader.Close()
	_ = part.logsWriter.Close()
	_ = part.profileReader.Close()
	_ = part.profileWriter.Close()
}
