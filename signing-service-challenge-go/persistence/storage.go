package persistence

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge-go/domain"
)

// Storage Interface - methods to manage the device data
type Storage interface {
	SaveDevice(device *domain.SignatureDevice) string
	GetDevice(deviceID string) (*domain.SignatureDevice, error)
	ListDevices() ([]*domain.SignatureDevice, error)
}
