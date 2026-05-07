package auth

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	Repo      *Repository
	Redis     *redis.Client
	JWTSecret string
}

// error definitions
var (
	ErrInvalidEmail       = errors.New("INVALID_EMAIL")
	ErrInvalidCredentials = errors.New("INVALID_CREDENTIALS")
	ErrUserInactive       = errors.New("USER_INACTIVE")
	ErrBatchInactive      = errors.New("BATCH_INACTIVE")
)

//
// ==========================
// PUBLIC METHODS (ENTRY POINT)
// ==========================
//

// LOGIN
func (s *Service) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {

	user, err := s.Repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrInvalidEmail
	}

	if err := s.validateLogin(user, req.Password); err != nil {
		return nil, err
	}

	accessToken, jti, err := s.generateAccessToken(user.ID, user.Role, req.Platform)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	refreshToken, _, err := s.generateRefreshToken(user.ID, user.Role, req.Platform)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := s.storeSession(ctx, user.ID, user.Role, req.Platform, jti, accessToken, refreshToken); err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// REFRESH
func (s *Service) Refresh(ctx context.Context, refreshToken string) (*LoginResponse, error) {

	userID, role, platform, err := s.validateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	accessToken, _, err := s.generateAccessToken(userID, role, platform)
	if err != nil {
		return nil, err
	}

	newRefreshToken, _, err := s.generateRefreshToken(userID, role, platform)
	if err != nil {
		return nil, err
	}

	// rotate refresh token
	key := buildRefreshKey(userID, platform)
	err = s.Redis.Set(ctx, key, newRefreshToken, 7*24*time.Hour).Err()
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

//
// ==========================
// INTERNAL FLOW (PRIVATE)
// ==========================
//

// VALIDASI LOGIN
func (s *Service) validateLogin(user *UserAuthEntity, password string) error {

	if user.Status != 1 {
		return ErrUserInactive
	}
	
	var valid bool

	if user.Password != "" {
		valid = verifyHash(user.Password, password)
	}

	if !valid && user.KodePemulihan != "" {
		if password == user.KodePemulihan {
			valid = true
		}
	}

	if user.Role == "peserta" && user.BatchStatus != 1 {
		return ErrBatchInactive
	}

	if !valid {
		return ErrInvalidCredentials
	}

	return nil
}

// SIMPAN SESSION
func (s *Service) storeSession(
	ctx context.Context,
	userID int,
	role string,
	platform string,
	jti string,
	accessToken string,
	refreshToken string,
) error {

	ttl := getAccessTTL(platform)

	sessionKey := s.buildSessionKey(userID, role, platform, jti)

	// single device (web & mobile)
	if platform == "web" || platform == "mobile" {
		s.Redis.Del(ctx, sessionKey)
	}

	err := s.Redis.Set(ctx, sessionKey, accessToken, ttl).Err()
	if err != nil {
		return err
	}

	refreshKey := buildRefreshKey(userID, platform)
	return s.Redis.Set(ctx, refreshKey, refreshToken, 7*24*time.Hour).Err()
}

// VALIDASI REFRESH TOKEN
func (s *Service) validateRefreshToken(refreshToken string) (int, string, string, error) {

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return 0, "", "", errors.New("invalid refresh token")
	}

	claims := token.Claims.(jwt.MapClaims)

	if claims["type"] != "refresh" {
		return 0, "", "", errors.New("invalid token type")
	}

	userIDStr := claims["sub"].(string)
	userID, _ := strconv.Atoi(userIDStr)

	role, _ := claims["role"].(string)
	platform := claims["platform"].(string)

	return userID, role, platform, nil
}

//
// ==========================
// TOKEN (JWT)
// ==========================
//

func (s *Service) generateAccessToken(userID int, role string, platform string) (string, string, error) {

	duration := getAccessTTL(platform)

	now := time.Now()
	jti := uuid.New().String()

	claims := jwt.MapClaims{
		"sub":      strconv.Itoa(userID),
		"role":     role,
		"platform": platform,
		"iss":      "ukaisyndrome",
		"aud":      "ukaisyndrome-users",
		"iat":      now.Unix(),
		"nbf":      now.Unix(),
		"exp":      now.Add(duration).Unix(),
		"jti":      jti,
		"type":     "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(s.JWTSecret))
	if err != nil {
		return "", "", err
	}

	return signed, jti, nil
}

func (s *Service) generateRefreshToken(userID int, role string, platform string) (string, string, error) {

	jti := uuid.New().String()

	claims := jwt.MapClaims{
		"sub":      strconv.Itoa(userID),
		"role":     role,
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