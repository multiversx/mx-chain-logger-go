package lifespan

import (
	"context"
	"fmt"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	logger "github.com/ElrondNetwork/elrond-go-logger"
)

const minEpochsLifeSpan = 1

type epochsLifeSpanner struct {
	*baseLifeSpanner
	spanInEpochs uint32
	cancelFunc   context.CancelFunc
}

func newEpochsLifeSpanner(es logger.EpochStartNotifier, epochsLifeSpan uint32) (*epochsLifeSpanner, error) {
	if check.IfNil(es) {
		return nil, fmt.Errorf("%w, epoch start notifier is nil", logger.ErrCreateEpochsLifeSpanner)
	}
	if epochsLifeSpan < minEpochsLifeSpan {
		return nil, fmt.Errorf("%w, min: %v, provided %v", logger.ErrCreateEpochsLifeSpanner, minEpochsLifeSpan, epochsLifeSpan)
	}

	els := &epochsLifeSpanner{
		spanInEpochs:    epochsLifeSpan,
		baseLifeSpanner: newBaseLifeSpanner(),
	}

	es.RegisterForEpochChangeConfirmed(
		func(epoch uint32) {
			if epoch%els.spanInEpochs == 0 {
				els.Notify(fmt.Sprintf("%v", epoch))
			}
		},
	)

	return els, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (sls *epochsLifeSpanner) IsInterfaceNil() bool {
	return sls == nil
}
