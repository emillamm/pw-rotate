package main

//import "fmt"
import "errors"

type Rotator interface {
	Rotate(user string, oldPw string, newPw string) error
	Ping(user string, pw string) error
}

var ErrNotAuthenticated = errors.New("failed to authenticate with provided credentials")
//type AuthenticationError struct {
//	Reason string
//}
//
//func (e AuthenticationError) Error() string {
//	return fmt.Sprintf("failed to authenticate with username and provided password: %s", e.Reason)
//}

