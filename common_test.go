package pwrotate

import (
	"testing"
)

// Test a Rotator
func RotatorTest(t *testing.T, r Rotator, user string, oldPw string, newPw string) {
	t.Helper()

	t.Run("authentication using oldPw should succeed and newPw should fail", func(t *testing.T) {
		if err := r.Ping(user, oldPw); err != nil {
			t.Errorf("authentication with oldPw %s not successful: %#v", oldPw, err)
		}
		if err := r.Ping(user, newPw); err == nil {
			t.Errorf("authentication with newPw %s was successful", newPw)
		}
	})

	t.Run("rotate should update oldPw to newPw", func(t *testing.T) {
		if err := r.Rotate(user, oldPw, newPw); err != nil {
			t.Errorf("failed to rotate passwords: %s", err)
		}
	})

	t.Run("authentication using oldPw should fail and newPw should succeed", func(t *testing.T) {
		if err := r.Ping(user, oldPw); err == nil {
			t.Errorf("authentication with oldPw %s was successful: %#v", oldPw, err)
		}
		if err := r.Ping(user, newPw); err != nil {
			t.Errorf("authentication with newPw %s not successful", newPw)
		}
	})

}

