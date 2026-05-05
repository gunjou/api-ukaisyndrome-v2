package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"

	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/scrypt"
)

// ==========================
// HASH VERIFICATION
// ==========================
func verifyHash(stored string, input string) bool {

	if stored == "" {
		return false
	}

	parts := strings.Split(stored, "$")
	if len(parts) < 3 {
		return false
	}

	prefix := parts[0]

	if strings.HasPrefix(prefix, "pbkdf2") {
		return verifyPBKDF2(stored, input)
	}

	if strings.HasPrefix(prefix, "scrypt") {
		return verifyScrypt(stored, input)
	}

	return false
}

func verifyPBKDF2(stored string, input string) bool {

	parts := strings.Split(stored, "$")
	if len(parts) != 3 {
		return false
	}

	meta := parts[0]
	salt := parts[1]
	hash := parts[2]

	metaParts := strings.Split(meta, ":")
	iter, _ := strconv.Atoi(metaParts[2])

	dk := pbkdf2.Key([]byte(input), []byte(salt), iter, 32, sha256.New)

	return hex.EncodeToString(dk) == hash
}

func verifyScrypt(stored string, input string) bool {

	parts := strings.Split(stored, "$")
	meta := parts[0]
	salt := parts[1]
	hash := parts[2]

	metaParts := strings.Split(meta, ":")
	N, _ := strconv.Atoi(metaParts[1])
	r, _ := strconv.Atoi(metaParts[2])
	p, _ := strconv.Atoi(metaParts[3])

	dk, err := scrypt.Key([]byte(input), []byte(salt), N, r, p, 32)
	if err != nil {
		return false
	}

	return hex.EncodeToString(dk) == hash
}