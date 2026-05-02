package module

import (
	"api-ukaisyndrome-v2/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Service *Service
}


// GetModul godoc
// @Summary Get modules for a user
// @Description Get a list of modules for the logged-in user based on their role
// @Tags Module
// @Accept json
// @Produce json
// @Security BearerAuth
// @Router /modul/peserta [get]
func (h *Handler) GetModul(c *fiber.Ctx) error {

	userID, ok := c.Locals("sub").(int)
	if !ok {
		return response.Error(c, 401, "unauthorized", "UNAUTHORIZED", nil)
	}

	role, _ := c.Locals("role").(string)

	data, err := h.Service.GetModul(c.Context(), userID, role)
	if err != nil {
		return response.Error(c, 500, err.Error(), "INTERNAL_ERROR", nil)
	}

	return response.Success(c, data)
}