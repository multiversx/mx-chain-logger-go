package pipes

import (
	"strings"
	"sync"

	logger "github.com/ElrondNetwork/elrond-go-logger"
)

type dummyLogsGatherer struct {
	lines []logger.LogLineHandler
	text  strings.Builder
	mutex sync.RWMutex
}

func (gatherer *dummyLogsGatherer) Write(p []byte) (n int, err error) {
	return 42, nil
}

func (gatherer *dummyLogsGatherer) Output(line logger.LogLineHandler) []byte {
	gatherer.mutex.Lock()
	defer gatherer.mutex.Unlock()

	gatherer.lines = append(gatherer.lines, line)
	gatherer.gatherText(line)
	return make([]byte, 0)
}

func (gatherer *dummyLogsGatherer) gatherText(line logger.LogLineHandler) {
	gatherer.text.WriteString(line.GetMessage() + "\n")

	for _, arg := range line.GetArgs() {
		gatherer.text.WriteString(arg + "\n")
	}
}

func (gatherer *dummyLogsGatherer) ContainsText(str string) bool {
	gatherer.mutex.RLock()
	defer gatherer.mutex.RUnlock()

	text := gatherer.text.String()
	return strings.Contains(text, str)
}

func (gatherer *dummyLogsGatherer) ContainsLogLine(loggerName string, level logger.LogLevel, message string) bool {
	gatherer.mutex.RLock()
	defer gatherer.mutex.RUnlock()

	for _, line := range gatherer.lines {
		matchedLevel := line.GetLogLevel() == int32(level)
		matchedMessage := line.GetMessage() == message
		matchedLogger := line.GetLoggerName() == loggerName

		if matchedLevel && matchedMessage && matchedLogger {
			return true
		}
	}

	return false
}

func (gatherer *dummyLogsGatherer) IsInterfaceNil() bool {
	return gatherer == nil
}
