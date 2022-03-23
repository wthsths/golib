package gl_aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

type aesEncryptor struct {
	key string
}

func NewAesEncryptor(key string) *aesEncryptor {
	return &aesEncryptor{
		key: key,
	}
}

func GenerateKey(bitSize int) (string, error) {
	if bitSize%8 != 0 || bitSize < 0 {
		return "", fmt.Errorf("invalid bit size")
	}

	randomKey := make([]byte, bitSize/8)
	_, err := rand.Read(randomKey)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(randomKey), nil
}

func (e *aesEncryptor) Encrypt(in string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(e.key))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(in), nil)
	return ciphertext, nil
}

func (e *aesEncryptor) Decrypt(in []byte) (string, error) {
	block, err := aes.NewCipher([]byte(e.key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	gcmNonceSize := gcm.NonceSize()
	if len(in) < gcmNonceSize {
		return "", fmt.Errorf("nonce size (%d) is greater than input size (%d)", gcmNonceSize, len(in))
	}

	nonce := in[:gcmNonceSize]
	in = in[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, in, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
