package file

import (
	"errors"
	"strings"
)

var (
	errDBIsClosed       = errors.New("DB is closed")
	errContextClosing   = errors.New("context closing")
	errInvalidParameter = errors.New("invalid parameter provided")
)

func isClosingError(err error) bool {
	if err == nil {
		return false
	}

	return strings.Contains(err.Error(), errDBIsClosed.Error()) ||
		strings.Contains(err.Error(), errContextClosing.Error())
}
