package token_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wengchaoxi/auth-proxy/pkg/token"
)

func TestTokenManager(t *testing.T) {
	tokenManager := token.NewTokenManager("passphrasewhichneedstobe32bytes!")

	t.Run("Generate token and parse Token", func(t *testing.T) {
		// Generate token
		expiresAt := time.Now().Add(time.Second).Unix()
		tokenStr := tokenManager.GenerateToken("id", expiresAt)

		// Verify token
		assert.Equal(t, true, tokenManager.IsValidToken(tokenStr))
		time.Sleep(2 * time.Second)
		assert.Equal(t, false, tokenManager.IsValidToken(tokenStr))

		// Parse token
		token, err := tokenManager.ParseTokenFromString(tokenStr)
		println(token.ExpiresAt)
		assert.NoError(t, err)
		assert.Equal(t, token.TokenId, "id")
	})
}
