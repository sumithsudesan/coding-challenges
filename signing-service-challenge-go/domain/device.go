package domain

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge-go/crypto"
)

// TODO: signature device domain model ...

// Interface - to support the signing operation
type KeyDataPair interface {
	Sign(data string) ([]byte, error)
}

// SignatureDevice - data model
type SignatureDevice struct {
	ID               string
	Label            string
	SignatureCounter int
	KeyPair          KeyDataPair
}

// Creates ECC key pair
func CreateECCKeyPair() (KeyDataPair, error) {
	generator := crypto.ECCGenerator{}
	pair, err := generator.Generate()
	if err != nil {
		return nil, err
	}
	return crypto.NewECDSAHelper(pair.Private), nil
}

// Creates Rsa key pair
func CreateRSAKeyPair() (KeyDataPair, error) {
	generator := crypto.RSAGenerator{}
	pair, err := generator.Generate()
	if err != nil {
		return nil, err
	}
	return crypto.NewRSAHelper(pair.Private), nil
}
