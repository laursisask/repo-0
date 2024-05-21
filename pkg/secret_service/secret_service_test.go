package secret_service


import (
	"testing"
)

func TestNewSecretFunc(t *testing.T) {
	// Test the NewSecret function
	secret := NewSecret("/org/freedesktop/secrets/1", "TestSecretString")
	if secret.Session != "/org/freedesktop/secrets/1" {
		t.Errorf("NewSecret function did not return expected value")
	}
	if string(secret.Value) != "TestSecretString" {
		t.Errorf("NewSecret function did not return expected value")
	}
}