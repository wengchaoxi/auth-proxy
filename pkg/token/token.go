package token

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/wengchaoxi/auth-proxy/pkg/crypto"
)

type Token struct {
	TokenId   string `json:"token_id"`
	ExpiresAt int64  `json:"expires_at"`
}

type TokenManager struct {
	token  Token
	crypto crypto.Crypto
}

func NewTokenManager(secretKey string) *TokenManager {
	c, _ := crypto.NewAESCrypto(&crypto.CryptoOptions{SecretKey: secretKey})
	return &TokenManager{crypto: c}
}

func (tm *TokenManager) GenerateToken(id string, expiresAt int64) string {
	// encrypt -> data
	str := fmt.Sprintf("%s\n%d", id, expiresAt)
	re, _ := tm.crypto.Encrypt(&crypto.CryptoOptions{Data: []byte(str)})
	return re
}

func (tm *TokenManager) ParseTokenFromString(x string) (Token, error) {
	// decrypt -> data
	str, err := tm.crypto.Decrypt(&crypto.CryptoOptions{Data: []byte(x)})
	if err != nil {
		return Token{}, err
	}
	strList := strings.Split(str, "\n")
	tm.token.TokenId = strList[0]
	tm.token.ExpiresAt, _ = strconv.ParseInt(strList[1], 10, 64)
	return tm.token, nil
}

func (tm *TokenManager) IsValidToken(x string) bool {
	_, err := tm.ParseTokenFromString(x)
	if err != nil {
		return false
	}
	if time.Now().Unix() > tm.token.ExpiresAt {
		return false
	}
	return true
}
