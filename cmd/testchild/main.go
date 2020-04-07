package main

import (
	"fmt"
	"os"

	"github.com/ElrondNetwork/elrond-go-logger/marshal"
	"github.com/ElrondNetwork/elrond-go-logger/pipes"
)

const (
	fileDescriptorProfileReader = 3
	fileDescriptorLogsWriter    = 4
)

func main() {
	errCode, errMessage := doMain()
	if errCode != 0 {
		fmt.Fprintln(os.Stderr, errMessage)
		os.Exit(errCode)
	}
}

// doMain returns (error code, error message)
func doMain() (int, string) {
	profileReader := getPipeFile(fileDescriptorProfileReader)
	if profileReader == nil {
		return 42, "Cannot get pipe file: [profileReader]"
	}

	logsWriter := getPipeFile(fileDescriptorLogsWriter)
	if logsWriter == nil {
		return 42, "Cannot get pipe file: [logsWriter]"
	}

	part, err := pipes.NewChildPart(profileReader, logsWriter, &marshal.JSONMarshalizer{})
	if err != nil {
		return 42, fmt.Sprintf("Cannot create ChildPart: %v", err)
	}

	err = part.StartLoop()
	if err != nil {
		return 42, fmt.Sprintf("Ended loop: %v", err)
	}

	return 0, ""
}

func getPipeFile(fileDescriptor uintptr) *os.File {
	file := os.NewFile(fileDescriptor, fmt.Sprintf("/proc/self/fd/%d", fileDescriptor))
	return file
}
