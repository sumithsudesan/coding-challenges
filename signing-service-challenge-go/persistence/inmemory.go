package persistence

import (
	"errors"
	"fmt"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge-go/domain"
)

// TODO: in-memory persistence ...

// InMemoryStorage - keep data in memory
// based Storage interface
type InMemoryStorage struct {
	devicesMap map[string]*domain.SignatureDevice
	mu         sync.RWMutex
}

// Create new InMemoryStorage instance
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		devicesMap: make(map[string]*domain.SignatureDevice),
	}
}

// SaveDevice - Saves the device details based on the device ID
func (s *InMemoryStorage) SaveDevice(device *domain.SignatureDevice) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if existingDevice, exists := s.devicesMap[device.ID]; exists {
		// Key already exists, return the existing ID
		return existingDevice.ID
	}
	s.devicesMap[device.ID] = device
	return device.ID
}

// GetDevice - Gets the device details based on the device ID
func (s *InMemoryStorage) GetDevice(deviceID string) (*domain.SignatureDevice, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	device, exists := s.devicesMap[deviceID]
	if !exists {
		errMsg := fmt.Sprintf("Device not found, device ID: %s", deviceID)
		return &domain.SignatureDevice{}, errors.New(errMsg)
	}
	return device, nil
}

// ListDevices - List all device details
func (s *InMemoryStorage) ListDevices() ([]*domain.SignatureDevice, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var devicesList []*domain.SignatureDevice
	for _, device := range s.devicesMap {
		devicesList = append(devicesList, device)
	}
	return devicesList, nil
}
