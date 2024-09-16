package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
)

type ECDSAHelper struct {
}

// Get the signature
func (h ECDSAHelper) Sign(dataToBeSigned []byte) ([]byte, error) {
	generator := ECCGenerator{}
	pair, err := generator.Generate()
	if err != nil {
		return nil, err
	}
	hash := sha256.New()
	hash.Write([]byte(dataToBeSigned))
	hashed := hash.Sum(nil)
	r, s, err := ecdsa.Sign(rand.Reader, pair.Private, hashed)
	if err != nil {
		return nil, err
	}
	return append(r.Bytes(), s.Bytes()...), nil
}
