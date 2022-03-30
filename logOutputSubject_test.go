package logger_test

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-go-logger/mock"
	"github.com/ElrondNetwork/elrond-go-logger/proto"
	"github.com/stretchr/testify/assert"
)

const testString1 = "DEBUG[2022-03-28 13:22:34.061] [consensus/spos/bls] [2/0/2/(END_ROUND)] step 3: block header final info has been received PubKeysBitmap = 1f AggregateSignature = 25f831bdb0801891a46b3b08a7bb11e306ad2e21d801a17312402a9d8bfc3ba76a4b97b42a8bc5ef533c471c47274c18 LeaderSignature = b2036b8db0bcaa7336e38f940b5f88706dc30afb6324693d01a93e9c47776ded31a195ac081b0c4274ed5c1354815292\n"
const testString2 = "DEBUG[2022-03-28 13:22:34.027] [..cess/coordinator] [2/0/2/(BLOCK)] elapsed time to processMiniBlocksToMe    time [s] = 90.747Âµs \n"

func TestNewLogOutputSubject(t *testing.T) {
	t.Parallel()

	los := logger.NewLogOutputSubject()
	assert.NotNil(t, los)
	assert.False(t, los.IsInterfaceNil())
}

// ------- AddObserver

func TestLogOutputSubject_AddObserverNilWriterShouldError(t *testing.T) {
	t.Parallel()

	los := logger.NewLogOutputSubject()

	err := los.AddObserver(nil, &mock.FormatterStub{})

	assert.Equal(t, logger.ErrNilWriter, err)
}

func TestLogOutputSubject_AddObserverNilFormatterShouldError(t *testing.T) {
	t.Parallel()

	los := logger.NewLogOutputSubject()

	err := los.AddObserver(&mock.WriterStub{}, nil)

	assert.Equal(t, logger.ErrNilFormatter, err)
}

func TestLogOutputSubject_AddObserverShouldWork(t *testing.T) {
	t.Parallel()

	los := logger.NewLogOutputSubject()

	err := los.AddObserver(&mock.WriterStub{}, &mock.FormatterStub{})
	writers, formatters := los.Observers()

	assert.Nil(t, err)
	assert.Equal(t, 1, len(writers))
	assert.Equal(t, 1, len(formatters))
}

// ------- Output

func TestLogOutputSubject_OutputNoObserversShouldDoNothing(t *testing.T) {
	t.Parallel()

	los := logger.NewLogOutputSubject()

	los.Output(nil)
}

func TestLogOutputSubject_OutputShouldCallFormatterAndWriter(t *testing.T) {
	t.Parallel()

	var formatterCalled = int32(0)
	var writerCalled = int32(0)
	los := logger.NewLogOutputSubject()
	_ = los.AddObserver(
		&mock.WriterStub{
			WriteCalled: func(p []byte) (n int, err error) {
				atomic.AddInt32(&writerCalled, 1)
				return 0, nil
			},
		},
		&mock.FormatterStub{
			OutputCalled: func(line logger.LogLineHandler) []byte {
				atomic.AddInt32(&formatterCalled, 1)
				return nil
			},
		},
	)

	los.Output(nil)

	assert.Equal(t, int32(1), atomic.LoadInt32(&writerCalled))
	assert.Equal(t, int32(1), atomic.LoadInt32(&formatterCalled))
}

func TestLogOutputSubject_OutputShouldProduceCorrectString(t *testing.T) {
	t.Parallel()

	los := logger.NewLogOutputSubject()
	var writtenData []byte
	_ = los.AddObserver(
		&mock.WriterStub{
			WriteCalled: func(p []byte) (n int, err error) {
				writtenData = p
				return 0, nil
			},
		},
		&logger.PlainFormatter{},
	)

	logLine := &logger.LogLine{
		LoggerName:  "",
		Correlation: proto.LogCorrelationMessage{},
		Message:     "message",
		LogLevel:    logger.LogDebug,
		Args: []interface{}{
			"int", 1,
			"ASCII string", "plain text \n",
			"non-ASCII string", "Âµs",
			"time.Duration", time.Microsecond*4 + time.Nanosecond,
			"byte slice", []byte("aaa"),
			"error", errors.New("an error"),
			"bool", true,
		},
		Timestamp: time.Date(2022, 03, 30, 15, 47, 52, 0, time.Local),
	}

	los.Output(logLine)

	expectedString := `DEBUG[2022-03-30 15:47:52.000]   message                                  int = 1 ASCII string = plain text 
 non-ASCII string = c382c2b573 time.Duration = 4.001µs byte slice = 616161 error = an error bool = true 
`

	assert.Equal(t, expectedString, string(writtenData))
}

func TestLogOutputSubject_OutputCalledConcurrentShouldWork(t *testing.T) {
	t.Parallel()

	var formatterCalled = int32(0)
	var writerCalled = int32(0)
	los := logger.NewLogOutputSubject()
	_ = los.AddObserver(
		&mock.WriterStub{
			WriteCalled: func(p []byte) (n int, err error) {
				atomic.AddInt32(&writerCalled, 1)
				return 0, nil
			},
		},
		&mock.FormatterStub{
			OutputCalled: func(line logger.LogLineHandler) []byte {
				atomic.AddInt32(&formatterCalled, 1)
				return nil
			},
		},
	)

	numCalls := 1000
	wg := &sync.WaitGroup{}
	wg.Add(numCalls)
	for i := 0; i < numCalls; i++ {
		go func() {
			time.Sleep(time.Millisecond)
			los.Output(nil)
			wg.Done()
		}()
	}

	wg.Wait()

	assert.Equal(t, int32(numCalls), atomic.LoadInt32(&writerCalled))
	assert.Equal(t, int32(numCalls), atomic.LoadInt32(&formatterCalled))
}

// ------- RemoveObserver

func TestLogOutputSubject_RemoveObserverNilWriterShouldError(t *testing.T) {
	t.Parallel()

	los := logger.NewLogOutputSubject()

	err := los.RemoveObserver(nil)

	assert.Equal(t, logger.ErrNilWriter, err)
}

func TestLogOutputSubject_RemoveObserverEmptyListShouldError(t *testing.T) {
	t.Parallel()

	los := logger.NewLogOutputSubject()

	err := los.RemoveObserver(&mock.WriterStub{})

	assert.Equal(t, logger.ErrWriterNotFound, err)
}

func TestLogOutputSubject_RemoveObserverWriterNotFoundShouldError(t *testing.T) {
	t.Parallel()

	los := logger.NewLogOutputSubject()
	_ = los.AddObserver(&mock.WriterStub{}, &mock.FormatterStub{})
	_ = los.AddObserver(&mock.WriterStub{}, &mock.FormatterStub{})

	err := los.RemoveObserver(&mock.WriterStub{})

	assert.Equal(t, logger.ErrWriterNotFound, err)
}

func TestLogOutputSubject_RemoveObserverOneElementShouldWork(t *testing.T) {
	t.Parallel()

	los := logger.NewLogOutputSubject()
	w := &mock.WriterStub{}
	_ = los.AddObserver(w, &mock.FormatterStub{})

	err := los.RemoveObserver(w)
	writers, formatters := los.Observers()

	assert.Nil(t, err)
	assert.Equal(t, 0, len(writers))
	assert.Equal(t, 0, len(formatters))
}

func TestLogOutputSubject_RemoveObserverLastElementShouldWork(t *testing.T) {
	t.Parallel()

	los := logger.NewLogOutputSubject()
	_ = los.AddObserver(&mock.WriterStub{}, &mock.FormatterStub{})
	w := &mock.WriterStub{}
	_ = los.AddObserver(w, &mock.FormatterStub{})

	err := los.RemoveObserver(w)
	writers, formatters := los.Observers()

	assert.Nil(t, err)
	assert.Equal(t, 1, len(writers))
	assert.Equal(t, 1, len(formatters))
}

func TestLogOutputSubject_RemoveObserverMiddleElementShouldWork(t *testing.T) {
	t.Parallel()

	los := logger.NewLogOutputSubject()
	_ = los.AddObserver(&mock.WriterStub{}, &mock.FormatterStub{})
	w := &mock.WriterStub{}
	_ = los.AddObserver(w, &mock.FormatterStub{})
	_ = los.AddObserver(&mock.WriterStub{}, &mock.FormatterStub{})

	err := los.RemoveObserver(w)
	writers, formatters := los.Observers()

	assert.Nil(t, err)
	assert.Equal(t, 2, len(writers))
	assert.Equal(t, 2, len(formatters))
}

func TestLogOutputSubject_ClearObservers(t *testing.T) {
	t.Parallel()

	los := logger.NewLogOutputSubject()

	w := &mock.WriterStub{}
	_ = los.AddObserver(w, &mock.FormatterStub{})

	obs, _ := los.Observers()
	assert.Equal(t, 1, len(obs))

	los.ClearObservers()

	obs, _ = los.Observers()
	assert.Equal(t, 0, len(obs))
}

func TestIsASCII(t *testing.T) {
	t.Parallel()

	assert.True(t, logger.IsASCII("ascii TEXT 1234 \\~&&\b\t\n"))
	assert.False(t, logger.IsASCII("µs"))
}

func BenchmarkIsASCII(b *testing.B) {
	b.Run("ASCII string", func(b *testing.B) {
		// should be < 150ns/op for the provided string
		for i := 0; i < b.N; i++ {
			_ = logger.IsASCII(testString1)
		}
	})
	b.Run("non ASCII string", func(b *testing.B) {
		// should be < 50ns/op for the provided string
		for i := 0; i < b.N; i++ {
			_ = logger.IsASCII(testString2)
		}
	})
}
