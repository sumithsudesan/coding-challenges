package persistence

import (
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge-go/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge-go/persistence"
)

// Test storage

func TestInMemoryStorage(t *testing.T) {
	store := persistence.NewInMemoryStorage()

	device := domain.SignatureDevice{
		ID:               "device123",
		SignatureCounter: 1,
		Label:            "My Device",
	}

	// Test CreateDevice
	store.SaveDevice(&device)

	// Test GetDevice
	dev, err := store.GetDevice("device123")
	if err != nil {
		t.Errorf("GetDevice() failed: %v", err)
	}

	if device != *dev {
		t.Errorf("GetDevice() returned unexpected device: got %v, want %+v", dev, &device)
	}

	// Test ListDevices
	lists, err := store.ListDevices()
	if err != nil {
		t.Errorf("ListDevices() failed: %v", err)
	}

	if len(lists) != 1 {
		t.Errorf("ListDevices() returned expected number of devices: got %d, want 1", len(lists))
	}
}
