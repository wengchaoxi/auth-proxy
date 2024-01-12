package crypto

type CryptoOptions struct {
	Data      []byte
	SecretKey string
}

type Crypto interface {
	Encrypt(*CryptoOptions) (string, error)
	Decrypt(*CryptoOptions) (string, error)
}
