package pkg

import (
	"testing"

	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestSetupDBEngine(t *testing.T) {
	SetupDBEngine()
	if DBEngine == nil || DBEngine.Ping() != nil {
		t.Errorf("Expected DBEngine to be initialized, got nil")
	}
	defer DBEngine.Close()
}

func TestSetupTestDBEngine(t *testing.T) {
	SetupTestDBEngine()
	if DBEngine == nil || DBEngine.Ping() != nil {
		t.Errorf("Expected DBEngine to be initialized, got nil")
	}
	defer DBEngine.Close()
}
