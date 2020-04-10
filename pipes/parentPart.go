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
	childName          string
	messenger          *ParentMessenger
	logLinesSink       logger.Logger
	textOutputSink     logger.Logger
	logLineMarshalizer logger.Marshalizer
	loopState          partLoopState

	logsReader    *os.File
	logsWriter    *os.File
	profileReader *os.File
	profileWriter *os.File
}

// NewParentPart creates a new logs receiver part (in the parent process)
func NewParentPart(childName string, logLineMarshalizer logger.Marshalizer) (*parentPart, error) {
	part := &parentPart{
		childName:          childName,
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

func (part *parentPart) StartLoop(childStdout io.Reader, childStderr io.Reader) error {
	if !part.loopState.isInit() {
		return ErrInvalidOperationGivenPartLoopState
	}

	logger.SubscribeToProfileChange(part)
	part.forwardProfile()
	part.continuouslyRead(childStdout, childStderr)
	part.loopState.setRunning()
	return nil
}

func (part *parentPart) OnProfileChanged() {
	part.forwardProfile()
}

func (part *parentPart) forwardProfile() {
	profile := logger.GetCurrentProfile()
	part.messenger.SendProfile(profile)
}

func (part *parentPart) continuouslyRead(childStdout io.Reader, childStderr io.Reader) {
	stdoutReader := bufio.NewReader(childStdout)
	stderrReader := bufio.NewReader(childStderr)

	go func() {
		for {
			if !part.loopState.isRunning() {
				break
			}

			logLine, err := part.messenger.ReadLogLine()
			if err != nil {
				part.logLinesSink.Error("continuouslyReadLogLines error", "err", err)
				break
			}

			part.logLinesSink.Log(logLine)
		}
	}()

	go func() {
		for {
			if !part.loopState.isRunning() {
				break
			}

			textLine, err := stdoutReader.ReadString('\n')
			if err != nil {
				break
			}

			textLine = strings.TrimSpace(textLine)
			part.textOutputSink.Trace(part.childName, "line", textLine)
		}
	}()

	go func() {
		for {
			if !part.loopState.isRunning() {
				break
			}

			textLine, err := stderrReader.ReadString('\n')
			if err != nil {
				break
			}

			textLine = strings.TrimSpace(textLine)
			part.textOutputSink.Error(part.childName, "line", textLine)
		}
	}()
}

func (part *parentPart) StopLoop() {
	part.loopState.setStopped()
	logger.UnsubscribeFromProfileChange(part)

	_ = part.logsReader.Close()
	_ = part.logsWriter.Close()
	_ = part.profileReader.Close()
	_ = part.profileWriter.Close()
}
