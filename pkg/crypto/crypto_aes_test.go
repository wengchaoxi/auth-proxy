package crypto_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wengchaoxi/auth-proxy/pkg/crypto"
)

func TestAESCrypto(t *testing.T) {
	data := []byte("hello world!")
	secretKey := "passphrasewhichneedstobe32bytes!"

	c, err := crypto.NewAESCrypto(&crypto.CryptoOptions{SecretKey: secretKey})
	assert.NoError(t, err)

	t.Run("encrypt and decrypt", func(t *testing.T) {
		ciphertext, _ := c.Encrypt(&crypto.CryptoOptions{Data: data})
		plaintext, _ := c.Decrypt(&crypto.CryptoOptions{Data: []byte(ciphertext)})
		fmt.Println("ciphertext:", ciphertext)
		fmt.Println("plaintext:", string(plaintext))
		assert.Equal(t, string(data), plaintext)
	})

	t.Run("ciphertext are not repeated", func(t *testing.T) {
		size := 10000
		set := make(map[string]bool, size)
		// encrypt
		for i := 0; i < size; i++ {
			ciphertext, _ := c.Encrypt(&crypto.CryptoOptions{Data: data})
			set[ciphertext] = true
		}
		assert.Equal(t, size, len(set))
		// decrypt
		for k := range set {
			plaintext, _ := c.Decrypt(&crypto.CryptoOptions{Data: []byte(k)})
			assert.Equal(t, string(data), plaintext)
		}
	})

}
