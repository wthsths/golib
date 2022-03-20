package gl_crypto

import (
	crand "crypto/rand"
	"encoding/binary"
	"log"
	"math/rand"
	"strings"
)

type passwordGenerator struct {
	rnd *rand.Rand
}

func NewPasswordGenerator() *passwordGenerator {
	var src cryptoSource
	newRnd := rand.New(src)
	newRnd.Seed(newRnd.Int63())

	return &passwordGenerator{rnd: newRnd}
}

func (g *passwordGenerator) GeneratePassword(length int) string {
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789" +
		"!?_-")
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[g.rnd.Intn(len(chars))])
	}
	return b.String()
}

type cryptoSource struct{}

func (s cryptoSource) Seed(seed int64) {}

func (s cryptoSource) Int63() int64 {
	return int64(s.Uint64() & ^uint64(1<<63))
}

func (s cryptoSource) Uint64() (v uint64) {
	err := binary.Read(crand.Reader, binary.BigEndian, &v)
	if err != nil {
		log.Fatal(err)
	}
	return v
}
