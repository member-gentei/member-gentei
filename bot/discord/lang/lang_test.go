package lang

import "testing"

func TestNewBundle(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("caught panic creating bundle: %+v", r)
		}
	}()
	NewBundle()
}
