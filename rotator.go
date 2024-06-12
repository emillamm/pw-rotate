package pwrotate

import "errors"

type Rotator interface {
	Rotate(user string, oldPw string, newPw string) error
	Ping(user string, pw string) error
}

var ErrNotAuthenticated = errors.New("failed to authenticate with provided credentials")

