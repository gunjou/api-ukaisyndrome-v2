package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/scrypt"
)

type Service struct {
	Repo      *Repository
	Redis     *redis.Client
	JWTSecret string
}

// ==========================
// LOGIN
// ==========================
func (s *Service) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {

	user, err := s.Repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if user.Status != 1 {
		return nil, errors.New("user inactive")
	}

	// ==========================
	// VALIDASI LOGIN
	// ==========================
	var valid bool

	// cek password
	if user.Password != "" {
		valid = verifyHash(user.Password, req.Password)
	}

	// fallback ke kode pemulihan
	if !valid && user.KodePemulihan != "" {
		if req.Password == user.KodePemulihan {
			valid = true
		}
	}

	if !valid {
		return nil, errors.New("invalid credentials")
	}

	// ==========================
	// TOKEN
	// ==========================
	accessToken, err := s.generateAccessToken(user.ID, user.Role, req.Platform)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	refreshToken, _, err := s.generateRefreshToken(user.ID, req.Platform)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// ==========================
	// REDIS SESSION
	// ==========================
	key := s.buildSessionKey(user.ID, user.Role, req.Platform)

	err = s.Redis.Set(ctx, key, refreshToken, 7*24*time.Hour).Err()
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// ==========================
// JWT
// ==========================
func (s *Service) generateAccessToken(userID int, role string, platform string) (string, error) {
	claims := jwt.MapClaims{
		"sub":      strconv.Itoa(userID),
		"role":     role,
		"platform": platform,
		"iss":      "ukaisyndrome",
		"aud":      "ukaisyndrome-users",
		"iat":      time.Now().Unix(),
		"nbf":      time.Now().Unix(),
		"exp":      time.Now().Add(90 * time.Minute).Unix(),
		"jti":      uuid.New().String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.JWTSecret))
}

func (s *Service) generateRefreshToken(userID int, platform string) (string, string, error) {

	jti := uuid.New().String()

	claims := jwt.MapClaims{
		"sub":      strconv.Itoa(userID),
		"platform": platform,
		"iss":      "ukaisyndrome",
		"aud":      "ukaisyndrome-auth",
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(7 * 24 * time.Hour).Unix(),
		"jti":      jti,
		"type":     "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(s.JWTSecret))
	if err != nil {
		return "", "", err
	}

	return signed, jti, nil
}

// ==========================
// SESSION KEY
// ==========================
func (s *Service) buildSessionKey(userID int, role string, platform string) string {
	return "session:" + role + ":" + platform + ":" + strconv.Itoa(userID)
}

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

	// 🔥 INI BAGIAN YANG KITA DEBUG
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