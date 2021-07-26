package lifespan

import (
	"context"
	"fmt"
	"time"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	logger "github.com/ElrondNetwork/elrond-go-logger"
)

const minMBLifeSpan = 1
const minRefreshInterval = time.Second

type sizeLifeSpanner struct {
	*baseLifeSpanner
	spanInMB        uint32
	refreshInterVal time.Duration
	cancelFunc      context.CancelFunc
	currentFile     string
	fileSizeChecker logger.FileSizeCheckHandler
}

func newSizeLifeSpanner(fileSizeChecker logger.FileSizeCheckHandler, sizeLifeSpanInMB uint32, refreshInterval time.Duration) (*sizeLifeSpanner, error) {
	if check.IfNil(fileSizeChecker) {
		return nil, fmt.Errorf("newSizeLifeSpanner %w, nil file size checker", logger.ErrCreateSizeLifeSpanner)
	}

	if sizeLifeSpanInMB < minMBLifeSpan {
		return nil, fmt.Errorf("newSizeLifeSpanner %w, provided size %v, min %v MB", logger.ErrCreateSizeLifeSpanner, sizeLifeSpanInMB, minMBLifeSpan)
	}

	if refreshInterval < minRefreshInterval {
		return nil, fmt.Errorf("newSizeLifeSpanner %w, provided refreshInterval %v, min %v", logger.ErrCreateSizeLifeSpanner, refreshInterval, minRefreshInterval)
	}

	sls := &sizeLifeSpanner{
		spanInMB:        sizeLifeSpanInMB * 1024 * 1024,
		baseLifeSpanner: newBaseLifeSpanner(),
		fileSizeChecker: fileSizeChecker,
	}

	return sls, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (sls *sizeLifeSpanner) IsInterfaceNil() bool {
	return sls == nil
}

// SetCurrentFile sets the file need for monitoring for the size
func (sls *sizeLifeSpanner) SetCurrentFile(path string) {
	if sls.cancelFunc != nil {
		sls.cancelFunc()
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	sls.cancelFunc = cancelFunc

	go sls.startTicker(ctx, path, int64(sls.spanInMB))
}

func (sls *sizeLifeSpanner) startTicker(ctx context.Context, path string, maxSize int64) {
	for {
		select {
		case <-time.After(sls.refreshInterVal):
			size, _ := sls.fileSizeChecker.GetSize(path)
			if size > maxSize {
				sls.lifeSpanChannel <- ""
			}
		case <-ctx.Done():
			return
		}
	}
}
