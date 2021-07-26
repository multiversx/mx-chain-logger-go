package mock

// EpochStartNotifierStub -
type EpochStartNotifierStub struct {
	epochFinalizedHandler                 []func(epoch uint32)
	NotifyEpochChangeConfirmedCalled      func(epoch uint32)
	RegisterForEpochChangeConfirmedCalled func(handler func(epoch uint32))
}

// RegisterForEpochChangeConfirmed -
func (esnm *EpochStartNotifierStub) RegisterForEpochChangeConfirmed(handler func(epoch uint32)) {
	if esnm.RegisterForEpochChangeConfirmedCalled != nil {
		esnm.RegisterForEpochChangeConfirmedCalled(handler)
	}

	esnm.epochFinalizedHandler = append(esnm.epochFinalizedHandler, handler)
}

// NotifyEpochChangeConfirmed -
func (esnm *EpochStartNotifierStub) NotifyEpochChangeConfirmed(epoch uint32) {
	if esnm.NotifyEpochChangeConfirmedCalled != nil {
		esnm.NotifyEpochChangeConfirmedCalled(epoch)
	}

	for _, hdl := range esnm.epochFinalizedHandler {
		hdl(epoch)
	}
}

// IsInterfaceNil checks if the underlying object is nil
func (esnm *EpochStartNotifierStub) IsInterfaceNil() bool {
	return esnm == nil
}
