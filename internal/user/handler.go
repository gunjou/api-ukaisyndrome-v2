package user

import (
	"api-ukaisyndrome-v2/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Service *Service
}


// Me godoc
// @Summary Get current user
// @Description Get logged in user profile based on JWT
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Router /user/me [get]
func (h *Handler) Me(c *fiber.Ctx) error {

	userID, ok := c.Locals("sub").(int)
	if !ok {
		return response.Error(c, 401, "unauthorized", "UNAUTHORIZED", nil)
	}

	role, ok := c.Locals("role").(string)
	if !ok {
		return response.Error(c, 401, "unauthorized", "UNAUTHORIZED", nil)
	}

	res, err := h.Service.GetMe(c.Context(), userID, role)
	if err != nil {
		return response.Error(c, 500, err.Error(), "INTERNAL_ERROR", nil)
	}

	return response.Success(c, res)
}