package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
)

type ECDSAHelper struct {
	privateKey *ecdsa.PrivateKey
}

// Creates new ECDSAHelper instnace
func NewECDSAHelper(privateKey *ecdsa.PrivateKey) ECDSAHelper {
	return ECDSAHelper{privateKey: privateKey}
}

// Get the signature
func (h ECDSAHelper) Sign(data string) ([]byte, error) {
	hash := sha256.New()
	hash.Write([]byte(data))
	hashed := hash.Sum(nil)
	r, s, err := ecdsa.Sign(rand.Reader, h.privateKey, hashed)
	if err != nil {
		return nil, err
	}
	return append(r.Bytes(), s.Bytes()...), nil
}
