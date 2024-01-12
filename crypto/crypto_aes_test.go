package crypto

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAESCrypto(t *testing.T) {
	data := []byte("hello world!")
	secretKey := "passphrasewhichneedstobe32bytes!"

	c, err := NewAESCrypto(&CryptoOptions{SecretKey: secretKey})
	assert.NoError(t, err)

	t.Run("encrypt and decrypt", func(t *testing.T) {
		ciphertext, _ := c.Encrypt(&CryptoOptions{Data: data})
		plaintext, _ := c.Decrypt(&CryptoOptions{Data: []byte(ciphertext)})
		fmt.Println("ciphertext:", ciphertext)
		fmt.Println("plaintext:", string(plaintext))
		assert.Equal(t, string(data), plaintext)
	})

	t.Run("ciphertext are not repeated", func(t *testing.T) {
		size := 10000
		set := make(map[string]bool, size)
		// encrypt
		for i := 0; i < size; i++ {
			ciphertext, _ := c.Encrypt(&CryptoOptions{Data: data})
			set[ciphertext] = true
		}
		assert.Equal(t, size, len(set))
		// decrypt
		for k := range set {
			plaintext, _ := c.Decrypt(&CryptoOptions{Data: []byte(k)})
			assert.Equal(t, string(data), plaintext)
		}
	})

}
