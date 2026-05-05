package shared

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"

	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/scrypt"
)

// ==========================
// HASH VERIFICATION HELPER
// ==========================
func VerifyHash(stored string, input string) bool {

	if stored == "" {
		return false
	}

	parts := strings.Split(stored, "$")
	if len(parts) < 3 {
		return false
	}

	prefix := parts[0]

	if strings.HasPrefix(prefix, "pbkdf2") {
		return VerifyPBKDF2(stored, input)
	}

	if strings.HasPrefix(prefix, "scrypt") {
		return VerifyScrypt(stored, input)
	}

	return false
}

func VerifyPBKDF2(stored string, input string) bool {

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

func VerifyScrypt(stored string, input string) bool {

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

func HashPassword(password string) (string, error) {

	salt := "random_salt" // ⚠️ nanti bisa improve
	iter := 1000000

	dk := pbkdf2.Key([]byte(password), []byte(salt), iter, 32, sha256.New)

	return "pbkdf2:sha256:" + strconv.Itoa(iter) + "$" + salt + "$" + hex.EncodeToString(dk), nil
}