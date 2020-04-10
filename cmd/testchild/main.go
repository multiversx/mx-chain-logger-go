package main

import (
	"fmt"
	goLog "log"
	"os"

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
}

func getPipeFile(fileDescriptor uintptr) *os.File {
	file := os.NewFile(fileDescriptor, fmt.Sprintf("/proc/self/fd/%d", fileDescriptor))
	return file
}
