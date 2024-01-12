package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

var (
	ErrSecretKeyFormat   = errors.New("err: the length of the secret key must be 32 bytes.")
	ErrDecryptDataFormat = errors.New("err: the format of the data to be decrypted is incorrect.")
)

type AESCrypto struct {
	secretKey string
}

func NewAESCrypto(co *CryptoOptions) (*AESCrypto, error) {
	if len(co.SecretKey) != 0x20 {
		return nil, ErrSecretKeyFormat
	}
	return &AESCrypto{
		secretKey: co.SecretKey,
	}, nil
}

// https://tutorialedge.net/golang/go-encrypt-decrypt-aes-tutorial/#our-encryption-client

func (a *AESCrypto) Encrypt(co *CryptoOptions) (string, error) {
	secretKey := a.secretKey

	// the length of the secret key must be 32 bytes
	if len(co.SecretKey) == 0x20 {
		secretKey = co.SecretKey
	} else if len(co.SecretKey) != 0 {
		return "", ErrSecretKeyFormat
	}

	c, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", err
	}

	// gcm or Galois/Counter Mode, is a mode of operation
	// for symmetric key cryptographic block ciphers
	// - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, co.Data, nil)
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// https://tutorialedge.net/golang/go-encrypt-decrypt-aes-tutorial/#our-decryption-client

func (a *AESCrypto) Decrypt(co *CryptoOptions) (string, error) {
	secretKey := a.secretKey

	if len(co.SecretKey) == 0x20 {
		secretKey = co.SecretKey
	} else if len(co.SecretKey) != 0 {
		return "", ErrSecretKeyFormat
	}

	c, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	data, _ := base64.URLEncoding.DecodeString(string(co.Data))
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", ErrDecryptDataFormat
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
