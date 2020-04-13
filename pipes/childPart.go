package pipes

import (
	"io"
	"os"

	logger "github.com/ElrondNetwork/elrond-go-logger"
)

var _ io.Writer = (*childPart)(nil)

var log = logger.GetOrCreate("pipes/childPart")

type childPart struct {
	messenger        *ChildMessenger
	outputSubject    logger.LogOutputHandler
	logLineFormatter logger.Formatter
	loopState        partLoopState
}

// NewChildPart creates a new logs sender part (in the child process)
func NewChildPart(
	profileReader *os.File,
	logsWriter *os.File,
	logLineMarshalizer logger.Marshalizer,
) (*childPart, error) {
	logLineFormatter, err := logger.NewLogLineWrapperFormatter(logLineMarshalizer)
	if err != nil {
		return nil, err
	}

	return &childPart{
		messenger:        NewChildMessenger(profileReader, logsWriter),
		outputSubject:    logger.GetLogOutputSubject(),
		logLineFormatter: logLineFormatter,
	}, nil
}

func (part *childPart) StartLoop() error {
	if !part.loopState.isInit() {
		return ErrInvalidOperationGivenPartLoopState
	}

	err := part.addAsObserver()
	if err != nil {
		return err
	}

	go part.continuouslyReadProfile()
	return nil
}

func (part *childPart) addAsObserver() error {
	part.outputSubject.ClearObservers()
	part.outputSubject.AddObserver(part, part.logLineFormatter)
	return nil
}

func (part *childPart) continuouslyReadProfile() {
	for {
		if !part.loopState.isRunning() {
			break
		}

		profile, err := part.messenger.ReadProfile()
		if err != nil {
			break
		}

		err = profile.Apply()
		log.Info("Profile change applied.")
	}
}

func (part *childPart) Write(logLineMarshalized []byte) (int, error) {
	return part.messenger.SendLogLine(logLineMarshalized)
}

func (part *childPart) StopLoop() {
	part.loopState.setStopped()
	part.outputSubject.RemoveObserver(part)
}
