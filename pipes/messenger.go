package pipes

import (
	"encoding/binary"
	"io"
	"os"
)

// Messenger intermediates communication (message exchange) via pipes
type Messenger struct {
	reader *os.File
	writer *os.File
}

// NewMessenger creates a new messenger
func NewMessenger(reader *os.File, writer *os.File) *Messenger {
	return &Messenger{
		reader: reader,
		writer: writer,
	}
}

// SendMessage sends a message over the pipe
func (messenger *Messenger) SendMessage(message []byte) (int, error) {
	length := len(message)
	err := messenger.sendMessageLength(length)
	if err != nil {
		return 0, err
	}

	_, err = messenger.writer.Write(message)
	if err != nil {
		return 0, err
	}

	return length, nil
}

func (messenger *Messenger) sendMessageLength(length int) error {
	buffer := make([]byte, sizeOfUint32)
	binary.LittleEndian.PutUint32(buffer, uint32(length))
	_, err := messenger.writer.Write(buffer)
	return err
}

// ReadMessage receives a message, reads it from the pipe
func (messenger *Messenger) ReadMessage() ([]byte, error) {
	length, err := messenger.readMessageLength()
	if err != nil {
		return nil, err
	}

	return messenger.readMessagePayload(length)
}

func (messenger *Messenger) readMessageLength() (int, error) {
	buffer := make([]byte, sizeOfUint32)
	_, err := io.ReadFull(messenger.reader, buffer)
	if err != nil {
		return 0, err
	}

	length := binary.LittleEndian.Uint32(buffer)
	return int(length), nil
}

func (messenger *Messenger) readMessagePayload(length int) ([]byte, error) {
	buffer := make([]byte, length)
	_, err := io.ReadFull(messenger.reader, buffer)
	if err != nil {
		return nil, err
	}

	return buffer, nil
}

// Close closes the pipes
func (messenger *Messenger) Close() {
	_ = messenger.reader.Close()
	_ = messenger.writer.Close()
}
