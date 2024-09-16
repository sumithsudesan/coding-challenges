package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

type RSASigner struct{}

func (r *RSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	generator := RSAGenerator{}
	pair, err := generator.Generate()
	if err != nil {
		return nil, err
	}

	hash := sha256.New()
	hash.Write([]byte(dataToBeSigned))
	hashed := hash.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, pair.Private, crypto.SHA256, hashed)
}
