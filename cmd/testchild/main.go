package main

import (
	"fmt"
	goLog "log"
	"os"
	"time"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-go-logger/marshal"
	"github.com/ElrondNetwork/elrond-go-logger/pipes"
)

const (
	fileDescriptorProfileReader = 3
	fileDescriptorLogsWriter    = 4
)

func main() {
	profileReader := getPipeFile(fileDescriptorProfileReader)
	if profileReader == nil {
		goLog.Fatal("Cannot get pipe file: [profileReader]")
	}

	logsWriter := getPipeFile(fileDescriptorLogsWriter)
	if logsWriter == nil {
		goLog.Fatal("Cannot get pipe file: [logsWriter]")
	}

	part := pipes.NewChildPart(profileReader, logsWriter, &marshal.JSONMarshalizer{})
	err := part.StartLoop()
	if err != nil {
		goLog.Fatal("Ended loop")
	}

	fooLog := logger.GetOrCreate("foo")
	barLog := logger.GetOrCreate("bar")

	fooLog.Info("foo-info")
	barLog.Info("bar-info")

	fooLog.Trace("foo-trace-no")
	barLog.Trace("bar-trace-no")

	go func() {
		fooLog.Info("foo-in-go")
		barLog.Info("bar-in-go")
	}()

	time.Sleep(3 * time.Second)

	fooLog.Trace("foo-trace-yes")
	barLog.Trace("bar-trace-yes")

	fmt.Println("Here's some stdout")
	fmt.Fprintln(os.Stderr, "Here's some stderr")

	time.Sleep(3 * time.Second)
}

func getPipeFile(fileDescriptor uintptr) *os.File {
	file := os.NewFile(fileDescriptor, fmt.Sprintf("/proc/self/fd/%d", fileDescriptor))
	return file
}
