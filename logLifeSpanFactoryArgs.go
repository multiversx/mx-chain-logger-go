package logger

// LogLifeSpanFactoryArgs contains the data needed for the creation of a logLifeSpanFactory
type LogLifeSpanFactoryArgs struct {
	EpochStartNotifierWithConfirm EpochStartNotifier
	LifeSpanType                  string
	RecreateEvery                 int
}
