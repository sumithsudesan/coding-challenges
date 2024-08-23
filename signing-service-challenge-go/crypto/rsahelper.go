package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

type RSAHelper struct {
	privateKey *rsa.PrivateKey
}

// Creates new RSAHelper instnace
func NewRSAHelper(privateKey *rsa.PrivateKey) RSAHelper {
	return RSAHelper{privateKey: privateKey}
}

// Get the signature
func (h RSAHelper) Sign(data string) ([]byte, error) {
	hash := sha256.New()
	hash.Write([]byte(data))
	hashed := hash.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, h.privateKey, crypto.SHA256, hashed)
}
