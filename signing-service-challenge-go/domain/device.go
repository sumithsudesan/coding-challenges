package domain

import (
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge-go/crypto"
)

// TODO: signature device domain model ...

// Interface - to support the signing operation
// SignatureDevice - data model
type SignatureDevice struct {
	ID               string
	Label            string
	SignatureCounter int
	Signer           crypto.Signer
	mu               sync.Mutex
}

// Increment counter
func (d *SignatureDevice) IncrementSignatureCounter() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.SignatureCounter++
}
