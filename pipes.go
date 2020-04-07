package logger

import (
	"encoding/binary"
	"io"
	"os"
	"time"
)

var _ io.Writer = (*pipeObserver)(nil)

const sizeOfUint32 = 4

type pipeObserver struct {
	writePipe *os.File
}

// NewPipeObserver creates a new observer that can be attached to any logger,
// and which writes the log data through a pipe.
// Ultimately, the data will be read by a "pipeObserverForwarder"
func NewPipeObserver(writePipe *os.File) *pipeObserver {
	return &pipeObserver{
		writePipe: writePipe,
	}
}

// Write sends a marshalized log line through the pipe, to be captured by the forwarder
// We have to ensure this is thread-safe
func (observer *pipeObserver) Write(logLineMarshalized []byte) (int, error) {
	length := len(logLineMarshalized)
	err := observer.writeLogLineLength(length)
	if err != nil {
		return 0, err
	}

	return observer.writePipe.Write(logLineMarshalized)
}

func (observer *pipeObserver) writeLogLineLength(length int) error {
	buffer := make([]byte, sizeOfUint32)
	binary.LittleEndian.PutUint32(buffer, uint32(length))
	_, err := observer.writePipe.Write(buffer)
	return err
}

type pipeObserverForwarder struct {
	readPipe    *os.File
	marshalizer Marshalizer
	loggerSink  *logger
}

// NewPipeObserverForwarder creates a new forwarder that reads log lines from a pipe
// and sends them to a generic logger sink
func NewPipeObserverForwarder(readPipe *os.File, marshalizer Marshalizer, loggerSink *logger) *pipeObserverForwarder {
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

func (forwarder *pipeObserverForwarder) readLogLine() (*LogLine, error) {
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

func (forwarder *pipeObserverForwarder) readLogLinePayload(length uint32) (*LogLineWrapper, error) {
	buffer := make([]byte, length)
	_, err := io.ReadFull(forwarder.readPipe, buffer)
	if err != nil {
		return nil, err
	}

	logLine := &LogLineWrapper{}
	err = forwarder.marshalizer.Unmarshal(logLine, buffer)
	if err != nil {
		return nil, err
	}

	return logLine, nil
}

func (forwarder *pipeObserverForwarder) recoverLogLine(wrapper *LogLineWrapper) *LogLine {
	logLine := &LogLine{
		LoggerName:  wrapper.LoggerName,
		Correlation: wrapper.Correlation,
		Message:     wrapper.Message,
		LogLevel:    LogLevel(wrapper.LogLevel),
		Args:        make([]interface{}, len(wrapper.Args)),
		Timestamp:   time.Unix(0, wrapper.Timestamp),
	}

	for i, str := range wrapper.Args {
		logLine.Args[i] = str
	}

	return logLine
}