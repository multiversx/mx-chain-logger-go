package logger

import (
	"encoding/json"
	"fmt"
	"sync"
)

var globalProfileChangeSubject *profileChangeSubject

func init() {
	globalProfileChangeSubject = NewProfileChangeSubject()
}

// Profile holds global logger options
type Profile struct {
	LogLevelPatterns string
	WithCorrelation  bool
	WithLoggerName   bool
}

// GetCurrentProfile gets the current logger profile
func GetCurrentProfile() Profile {
	return Profile{
		LogLevelPatterns: GetLogLevelPattern(),
		WithCorrelation:  IsEnabledCorrelation(),
		WithLoggerName:   IsEnabledLoggerName(),
	}
}

// UnmarshalProfile deserializes into a Profile object
func UnmarshalProfile(data []byte) (Profile, error) {
	profile := Profile{}
	err := json.Unmarshal(data, &profile)
	if err != nil {
		return Profile{}, err
	}

	return profile, nil
}

// Marshal serializes the Profile object
func (profile *Profile) Marshal() ([]byte, error) {
	data, err := json.Marshal(profile)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Apply sets the global logger options
func (profile *Profile) Apply() error {
	err := SetLogLevel(profile.LogLevelPatterns)
	if err != nil {
		return err
	}

	ToggleCorrelation(profile.WithCorrelation)
	ToggleLoggerName(profile.WithLoggerName)
	globalProfileChangeSubject.NotifyAll()

	return nil
}

func (profile *Profile) String() string {
	return fmt.Sprintf("[pattern=%s, with correlation=%t, with logger name=%t]",
		profile.LogLevelPatterns,
		profile.WithCorrelation,
		profile.WithLoggerName,
	)
}

type profileChangeSubject struct {
	observers []ProfileChangeObserver
	mutex     sync.RWMutex
}

// NewProfileChangeSubject -
func NewProfileChangeSubject() *profileChangeSubject {
	return &profileChangeSubject{
		observers: make([]ProfileChangeObserver, 0),
	}
}

// Subscribe -
func (subject *profileChangeSubject) Subscribe(observer ProfileChangeObserver) {
	subject.mutex.Lock()
	subject.observers = append(subject.observers)
	subject.mutex.Unlock()
}

// Unsubscribe -
func (subject *profileChangeSubject) Unsubscribe(observer ProfileChangeObserver) {
	subject.mutex.Lock()
	defer subject.mutex.Unlock()

	for i := 0; i < len(subject.observers); i++ {
		if subject.observers[i] == observer {
			subject.observers = append(subject.observers[0:i], subject.observers[i+1:]...)
		}
	}
}

func (subject *profileChangeSubject) NotifyAll() {
	subject.mutex.RLock()
	defer subject.mutex.RUnlock()

	for _, observer := range subject.observers {
		observer.OnProfileChanged()
	}
}
