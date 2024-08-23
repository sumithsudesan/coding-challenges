package persistence

import (
	"log"
)

// NewStorage - create storage based on type
func NewStorage(storageType string) Storage {
	switch storageType {
	case "in-memory":
		return NewInMemoryStorage()

	// Current, we have in -memory
	// in future we can add more storage types
	default:
		log.Fatalf("Unsupported storage type: %s", storageType)
		return nil
	}
}
