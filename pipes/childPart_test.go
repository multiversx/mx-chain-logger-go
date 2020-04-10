package pipes

import (
	"os"
	"sync"
	"testing"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-go-logger/marshal"
)

func TestChildPart_NoPanicWhenNoParent(t *testing.T) {
	// Bad pipes (no parent)
	profileReader := os.NewFile(4242, "/proc/self/fd/4242")
	logsWriter := os.NewFile(4343, "/proc/self/fd/4343")

	logLineMarshalizer := &marshal.JSONMarshalizer{}
	childLogger := logger.GetOrCreate("child-log")
	childPart := NewChildPart(profileReader, logsWriter, logLineMarshalizer)
	childPart.StartLoop()

	childLogger.Debug("foo")
	childLogger.Trace("bar")
}

func TestChildPart_ConcurrentWriteLogs(t *testing.T) {
	profileReader := os.NewFile(4242, "/proc/self/fd/4242")
	logsWriter := os.NewFile(4343, "/proc/self/fd/4343")

	logLineMarshalizer := &marshal.JSONMarshalizer{}
	childLogger := logger.GetOrCreate("child-log")
	childPart := NewChildPart(profileReader, logsWriter, logLineMarshalizer)
	childPart.StartLoop()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		for i := 0; i < 1000; i++ {
			childLogger.Debug("foo")
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < 1000; i++ {
			childLogger.Trace("bar")
		}
		wg.Done()
	}()

	wg.Wait()
}
