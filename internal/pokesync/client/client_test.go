package client

import "testing"

func TestBuildNumber_IsUpToDateWith(t *testing.T) {
	b1 := BuildNumber(17)
	b2 := BuildNumber(16)

	if !b1.IsUpToDateWith(b2) {
		t.Error("expected BuildNumber to be up-to-date")
	}
}

func TestBuildNumber_IsUpToDateWith2(t *testing.T) {
	b1 := BuildNumber(16)
	b2 := BuildNumber(17)

	if b1.IsUpToDateWith(b2) {
		t.Error("expected BuildNumber to not be up-to-date")
	}
}
