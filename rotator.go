package pwrotate

import "errors"

type Rotator interface {
	Rotate(user string, oldPw string, newPw string) error
	Ping(user string, pw string) error
}

var ErrAlreadyRotated = errors.New("password already rotated")

