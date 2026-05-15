package cdn

import (
	"api-ukaisyndrome-v2/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Service *Service
}

// =================================================
// GET MENTOR IMAGES
// =================================================

// GetMentorImages godoc
// @Summary Get mentor images
// @Description Get list mentor images from storage
// @Tags CDN
// @Accept json
// @Produce json
// @Router /cdn/mentor [get]
func (h *Handler) GetMentorImages(c *fiber.Ctx) error {

	data, err := h.Service.GetMentorImages()
	if err != nil {
		return response.Error(
			c,
			500,
			err.Error(),
			"INTERNAL_ERROR",
			nil,
		)
	}

	return response.Success(c, data)
}

// =================================================
// GET MODUL IMAGES
// =================================================

// GetModulImages godoc
// @Summary Get modul images
// @Description Get list modul images from storage
// @Tags CDN
// @Accept json
// @Produce json
// @Router /cdn/modul [get]
func (h *Handler) GetModulImages(c *fiber.Ctx) error {

	data, err := h.Service.GetModulImages()
	if err != nil {
		return response.Error(
			c,
			500,
			err.Error(),
			"INTERNAL_ERROR",
			nil,
		)
	}

	return response.Success(c, data)
}

// =================================================
// GET ADS
// =================================================

// GetAds godoc
// @Summary Get ads
// @Description Get list ads images and links
// @Tags CDN
// @Accept json
// @Produce json
// @Router /cdn/ads [get]
func (h *Handler) GetAds(c *fiber.Ctx) error {

	data, err := h.Service.GetAds()
	if err != nil {
		return response.Error(
			c,
			500,
			err.Error(),
			"INTERNAL_ERROR",
			nil,
		)
	}

	return response.Success(c, data)
}