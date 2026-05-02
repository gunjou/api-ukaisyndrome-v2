package auth

import (
	"api-ukaisyndrome-v2/pkg/response"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Service *Service
}

func RegisterRoutes(r fiber.Router, handler *Handler) {
	auth := r.Group("/auth")

	auth.Post("/login", handler.Login)
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
		return response.Error(c, 401, "invalid credentials", "AUTH_INVALID", nil)
	}

	_ = time.Since(start) // nanti bisa dipakai logging

	return response.Success(c, res)
}