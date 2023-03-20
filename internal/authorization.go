package togo

import (
	"encoding/hex"

	"golang.org/x/crypto/argon2"
)

type Argon2Params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	keyLength   uint32
}

type BasicAuthConfig struct {
	argonConfig *Argon2Params
	saltLength  int
}

func NewBasicAuthConfig() *BasicAuthConfig {
	return &BasicAuthConfig{
		argonConfig: &Argon2Params{
			memory:      64 * 1024, // 64MB
			iterations:  1,
			parallelism: 2,
			keyLength:   32,
		},
		saltLength: 16,
	}
}

func (c BasicAuthConfig) HashPassword(password string) (key, salt string) {
	salt = generateRandomString(c.saltLength)
	key = hex.EncodeToString(argon2.IDKey(
		[]byte(password),
		[]byte(salt),
		c.argonConfig.iterations,
		c.argonConfig.memory,
		c.argonConfig.parallelism,
		c.argonConfig.keyLength,
	))

	return
}

func (c BasicAuthConfig) VerifyPassword(password, key, salt string) bool {
	return hex.EncodeToString(argon2.IDKey(
		[]byte(password),
		[]byte(salt),
		c.argonConfig.iterations,
		c.argonConfig.memory,
		c.argonConfig.parallelism,
		c.argonConfig.keyLength,
	)) == key
}
