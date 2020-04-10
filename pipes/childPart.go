package pipes

import (
	"io"
	"os"

	logger "github.com/ElrondNetwork/elrond-go-logger"
)

var _ io.Writer = (*childPart)(nil)

var log = logger.GetOrCreate("pipes/childPart")

type childPart struct {
	messenger          *ChildMessenger
	outputSubject      logger.LogOutputHandler
	logLineMarshalizer logger.Marshalizer
}

// NewChildPart creates a new logs sender part (in the child process)
func NewChildPart(
	profileReader *os.File,
	logsWriter *os.File,
	logLineMarshalizer logger.Marshalizer,
) *childPart {
	messenger := NewChildMessenger(profileReader, logsWriter)
	outputSubject := logger.GetLogOutputSubject()

	return &childPart{
		messenger:          messenger,
		outputSubject:      outputSubject,
		logLineMarshalizer: logLineMarshalizer,
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
	logLineFormatter, err := logger.NewLogLineWrapperFormatter(part.logLineMarshalizer)
	if err != nil {
		return err
	}

	part.outputSubject.ClearObservers()
	part.outputSubject.AddObserver(part, logLineFormatter)
	return nil
}

func (part *childPart) continuouslyReadProfile() {
	for {
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

func (part *childPart) Close() {
}
