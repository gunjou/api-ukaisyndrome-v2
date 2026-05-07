package auth

import (
	"api-ukaisyndrome-v2/pkg/response"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Service *Service
}

func RegisterRoutes(r fiber.Router, handler *Handler) {
	auth := r.Group("/auth")

	auth.Post("/login", handler.Login)
	auth.Post("/refresh", handler.Refresh)
}

// Login godoc
// @Summary Login user
// @Description Login dengan email dan password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login Request"
// @Router /auth/login [post]
func (h *Handler) Login(c *fiber.Ctx) error {
	start := time.Now()

	var req LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, 400, "invalid request body", "AUTH_INVALID_REQUEST", nil)
	}

	ctx := c.UserContext()

	res, err := h.Service.Login(ctx, req)
	if err != nil {

		if errors.Is(err, ErrInvalidEmail) {
			return response.Error(c, 400, "invalid email", "AUTH_INVALID_EMAIL", nil)
		}

		if errors.Is(err, ErrUserInactive) {
			return response.Error(c, 403, "user inactive", "AUTH_USER_INACTIVE", nil)
		}

		if errors.Is(err, ErrInvalidCredentials) {
			return response.Error(c, 401, "invalid credentials", "AUTH_INVALID_CREDENTIALS", nil)
		}

		if errors.Is(err, ErrBatchInactive) {
			return response.Error(c, 403, "batch inactive", "AUTH_BATCH_INACTIVE", nil)
		}

		return response.Error(c, 500, "internal server error", "INTERNAL_ERROR", nil)
	}

	_ = time.Since(start) // nanti bisa dipakai logging

	return response.Success(c, res)
}


// Refresh godoc
// @Summary Refresh access token
// @Description Generate new access token using refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RefreshRequest true "Refresh Token"
// @Router /auth/refresh [post]
func (h *Handler) Refresh(c *fiber.Ctx) error {

	var req RefreshRequest

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, 400, "invalid request", "BAD_REQUEST", nil)
	}

	res, err := h.Service.Refresh(c.Context(), req.RefreshToken)
	if err != nil {
		return response.Error(c, 401, err.Error(), "UNAUTHORIZED", nil)
	}

	return response.Success(c, res)
}