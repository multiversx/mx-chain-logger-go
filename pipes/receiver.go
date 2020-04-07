package pipes

import (
	"encoding/binary"
	"io"
	"os"
)

// Receiver intermediates communication (message receiving) via pipes
type Receiver struct {
	reader *os.File
}

// NewReceiver creates a new receiver
func NewReceiver(reader *os.File) *Receiver {
	return &Receiver{
		reader: reader,
	}
}

// Receive receives a message, reads it from the pipe
func (receiver *Receiver) Receive() ([]byte, error) {
	length, err := receiver.receiveMessageLength()
	if err != nil {
		return nil, err
	}

	return receiver.readMessage(length)
}

func (receiver *Receiver) receiveMessageLength() (int, error) {
	buffer := make([]byte, sizeOfUint32)
	_, err := io.ReadFull(receiver.reader, buffer)
	if err != nil {
		return 0, err
	}

	length := binary.LittleEndian.Uint32(buffer)
	return int(length), nil
}

func (receiver *Receiver) readMessage(length int) ([]byte, error) {
	buffer := make([]byte, length)
	_, err := io.ReadFull(receiver.reader, buffer)
	if err != nil {
		return nil, err
	}

	return buffer, nil
}

// Shutdown closes the pipe
func (receiver *Receiver) Shutdown() error {
	err := receiver.reader.Close()
	return err
}
