package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"apistarter/internal/config"
)

func AESGCMEncrypt(plaintext, key, additionData []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("error creating the cipher block :%w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("error creating the gcm block :%w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("error generating the nonce array: %w", err)
	}

	cipherText := gcm.Seal(nonce, nonce, plaintext, additionData)
	return cipherText, nil
}

func AESGCMDecript(cipherText, key, additionData []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("error creating the cipher block :%w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("error creating the gcm block :%w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		return nil, fmt.Errorf("cipher text to short")
	}

	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]
	plainText, err := gcm.Open(nil, nonce, cipherText, additionData)
	if err != nil {
		return nil, fmt.Errorf("error decripting the data: %w", err)
	}
	return plainText, nil
}

type DataEncrypter struct {
	apiKey []byte
}

func NewDataEncripter(cfg *config.Configuration) *DataEncrypter {
	return &DataEncrypter{
		apiKey: []byte(cfg.Key[:16]),
	}
}

func (d *DataEncrypter) Encrypt(text, additionalData []byte) ([]byte, error) {
	return AESGCMEncrypt(text, d.apiKey, additionalData)
}

func (d *DataEncrypter) Decrypt(cipherText, additionalData []byte) ([]byte, error) {
	return AESGCMDecript(cipherText, d.apiKey, additionalData)
}

func (d *DataEncrypter) DecryptBase64(cipherText string, additionalData []byte) ([]byte, error) {
	cipherBytes, err := base64.RawStdEncoding.DecodeString(cipherText)
	if err != nil {
		return nil, fmt.Errorf("error parsing from base64: %w", err)
	}
	return AESGCMDecript(cipherBytes, d.apiKey, additionalData)
}
