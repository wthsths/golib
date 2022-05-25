package gl_aes

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Encrypt_Decrypt(t *testing.T) {
	key, err := GenerateKey(256)

	fmt.Println(key)
	fmt.Println(len([]byte(key)))

	if err != nil {
		t.Fatal(err)
	}

	encryptor := NewAesEncryptor(key)

	in := "test"

	out, err := encryptor.Encrypt(in)
	if err != nil {
		t.Fatal(err)
	}

	decryptedOut, err := encryptor.Decrypt(out)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, in, decryptedOut)
}
