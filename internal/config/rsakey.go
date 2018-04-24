package config

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

// ReadRSAKey will parse the supplied path and return an RSA Key
func ReadRSAKey(path string) (*rsa.PrivateKey, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Unable to open private key file: %v", err)
	}
	block, _ := pem.Decode(b)
	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse private key: %v", err)
	}
	return privKey, nil
}

// Encrypt a Base64 encoded message with the supplied public key
func Encrypt(pubKey *rsa.PublicKey, message string) (string, error) {
	encrypted, err := rsa.EncryptOAEP(
		sha1.New(),
		rand.Reader,
		pubKey,
		[]byte(message),
		[]byte("todo"),
	)
	if err != nil {
		return "", err
	}
	strEnc := base64.StdEncoding.EncodeToString(encrypted)
	return strEnc, nil
}

// Decrypt a Base64 encoded message with the supplied private key
func Decrypt(privKey *rsa.PrivateKey, message string) (string, error) {
	bytesEnc, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return "", fmt.Errorf("Unbable to decode base64 string: %v", err)
	}
	decrypted, err := rsa.DecryptOAEP(
		sha1.New(),
		rand.Reader,
		privKey,
		bytesEnc,
		[]byte("todo"),
	)
	if err != nil {
		return "", fmt.Errorf("Unable to rsa decrypt: %v", err)
	}
	return string(decrypted), nil
}
