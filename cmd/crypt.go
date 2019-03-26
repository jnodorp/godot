package cmd

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"

	log "github.com/sirupsen/logrus"
)

func encrypt(plaintext, password []byte) (string, error) {
	c, err := aes.NewCipher(hash(password))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.WithError(err).Fatal("failed to generate nonce")
	}

	sealed := gcm.Seal(nonce, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(sealed), nil
}

func decrypt(ciphertext, password []byte) (string, error) {
	sealed, err := base64.StdEncoding.DecodeString(string(ciphertext))
	if err != nil {
		return "", err
	}

	c, err := aes.NewCipher(hash(password))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(sealed) < nonceSize {
		return "", err
	}

	nonce, sealed := sealed[:nonceSize], sealed[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, sealed, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func hash(password []byte) []byte {
	bytes := sha256.Sum256([]byte(password))
	r := make([]byte, sha256.Size, sha256.Size)
	for i, b := range bytes {
		r[i] = b
	}

	return r
}
