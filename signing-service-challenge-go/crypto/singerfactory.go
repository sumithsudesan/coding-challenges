package crypto

import "errors"

func GetSigner(algorithm string) (Signer, error) {
	switch algorithm {
	case "RSA":
		return &RSASigner{}, nil
	case "ECC":
		return &ECDSAHelper{}, nil
	default:
		return nil, errors.New("unsupported algorithm")
	}
}
