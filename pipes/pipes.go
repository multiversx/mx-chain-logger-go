package pipes

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"time"

	logger "github.com/ElrondNetwork/elrond-go-logger"
)

const sizeOfUint32 = 4

type pipeObserverForwarder struct {
	readPipe    *os.File
	marshalizer logger.Marshalizer
	loggerSink  logger.Logger
}

// NewPipeObserverForwarder creates a new forwarder that reads log lines from a pipe
// and sends them to a generic logger sink
func NewPipeObserverForwarder(readPipe *os.File, marshalizer logger.Marshalizer, loggerSink logger.Logger) *pipeObserverForwarder {
	return &pipeObserverForwarder{
		readPipe:    readPipe,
		marshalizer: marshalizer,
		loggerSink:  loggerSink,
	}
}

func (forwarder *pipeObserverForwarder) StartFowarding() {
	go forwarder.continuouslyReadLogLines()
}

func (forwarder *pipeObserverForwarder) continuouslyReadLogLines() {
	for {
		logLine, err := forwarder.readLogLine()
		if err != nil {
			forwarder.loggerSink.Error("continuouslyReadLogLines error", "err", err)
			break
		}

		forwarder.loggerSink.Log(logLine)
	}
}

func (forwarder *pipeObserverForwarder) readLogLine() (*logger.LogLine, error) {
	length, err := forwarder.readLogLineLength()
	if err != nil {
		return nil, err
	}

	logLineWrapper, err := forwarder.readLogLinePayload(length)
	if err != nil {
		return nil, err
	}

	logLine := forwarder.recoverLogLine(logLineWrapper)
	return logLine, nil
}

func (forwarder *pipeObserverForwarder) readLogLineLength() (uint32, error) {
	buffer := make([]byte, sizeOfUint32)
	_, err := io.ReadFull(forwarder.readPipe, buffer)
	if err != nil {
		return 0, err
	}

	length := binary.LittleEndian.Uint32(buffer)
	return length, nil
}

func (forwarder *pipeObserverForwarder) readLogLinePayload(length uint32) (*logger.LogLineWrapper, error) {
	buffer := make([]byte, length)
	_, err := io.ReadFull(forwarder.readPipe, buffer)
	if err != nil {
		return nil, err
	}

	logLine := &logger.LogLineWrapper{}
	err = forwarder.marshalizer.Unmarshal(logLine, buffer)
	if err != nil {
		return nil, err
	}

	return logLine, nil
}

func (forwarder *pipeObserverForwarder) recoverLogLine(wrapper *logger.LogLineWrapper) *logger.LogLine {
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

type pipeProfileForwarder struct {
	writePipe *os.File
}

// NewPipeProfileForwarder creates a new profile forwarder,
// which forwards logging profiles through pipe
func NewPipeProfileForwarder(writePipe *os.File) *pipeProfileForwarder {
	return &pipeProfileForwarder{
		writePipe: writePipe,
	}
}

func (forwarder *pipeProfileForwarder) StartFowarding() {
	logger.SubscribeToProfileChange(forwarder)
	forwarder.forwardProfile()
}

func (forwarder *pipeProfileForwarder) OnProfileChanged() {
	forwarder.forwardProfile()
}

func (forwarder *pipeProfileForwarder) forwardProfile() {
	profile := logger.GetCurrentProfile()
	fmt.Println(profile)
}

func (forwarder *pipeProfileForwarder) Close() {
	logger.UnsubscribeFromProfileChange(forwarder)
}

type pipeProfileReceiver struct {
}

// TODO Messenger = sender + receiver.
