package logger_test

import (
	"fmt"
	"testing"

	"github.com/multiversx/mx-chain-core-go/marshal"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-logger-go/proto"
	"github.com/stretchr/testify/assert"
)

func TestLogLineWrapper_MarshalUnmarshalShouldWork(t *testing.T) {
	llw := generateLogLineWrapper()

	testMarshalUnmarshal(t, "gogo protobuf", &marshal.GogoProtoMarshalizer{}, llw)
}

func generateLogLineWrapper() logger.LogLineWrapper {
	return logger.LogLineWrapper{
		LogLineMessage: proto.LogLineMessage{
			Message:   "test message",
			LogLevel:  4,
			Args:      []string{"arg1", "arg2", "arg3", "arg4"},
			Timestamp: 11223344,
		},
	}
}

func testMarshalUnmarshal(t *testing.T, marshName string, marsh logger.Marshalizer, llw logger.LogLineWrapper) {
	llwCopyForAssert := llw

	buff, err := marsh.Marshal(&llw)
	assert.Nil(t, err)

	llwRecovered := &logger.LogLineWrapper{}
	err = marsh.Unmarshal(llwRecovered, buff)
	assert.Nil(t, err)

	assert.Equal(t, &llwCopyForAssert, llwRecovered, fmt.Sprintf("for marshalizer %v", marshName))
}
