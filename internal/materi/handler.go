package materi

import (
	"api-ukaisyndrome-v2/pkg/response"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Service *Service
}


// GetMateriPeserta godoc
// @Summary Get materi peserta
// @Description Get list materi berdasarkan modul
// @Tags Materi
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id_modul path int true "ID Modul"
// @Param type query string false "Filter tipe materi (video | document)"
// @Router /materi/peserta/{id_modul} [get]
func (h *Handler) GetMateriPeserta(c *fiber.Ctx) error {

	userID, ok := c.Locals("sub").(int)
	if !ok {
		return response.Error(c, 401, "unauthorized", "UNAUTHORIZED", nil)
	}

	modulID, err := c.ParamsInt("id_modul")
	if err != nil {
		return response.Error(c, 400, "invalid modul id", "BAD_REQUEST", nil)
	}

	// query param
	materiType := c.Query("type")

	var typePtr *string
	if materiType != "" {
		typePtr = &materiType
	}

	data, err := h.Service.GetMateriPeserta(c.Context(), userID, modulID, typePtr)
	if err != nil {
		return response.Error(c, 500, err.Error(), "INTERNAL_ERROR", nil)
	}

	return response.Success(c, data)
}

// GetMateriPrivatePeserta godoc
// @Summary Get materi private peserta
// @Description Get list materi private berdasarkan modul
// @Tags Materi
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param type query string false "Filter tipe materi (video | document)"
// @Router /materi/peserta/private [get]
func (h *Handler) GetMateriPrivatePeserta(c *fiber.Ctx) error {

	userID := c.Locals("sub").(int)

	// optional query
	materiType := c.Query("type")

	var typePtr *string
	if materiType != "" {
		typePtr = &materiType
	}

	data, err := h.Service.GetMateriPrivatePeserta(
		c.Context(),
		userID,
		typePtr,
	)

	if err != nil {
		return response.Error(c, 400, err.Error(), "BAD_REQUEST", nil)
	}

	return response.Success(c, data)
}


// OpenMateri godoc
// @Summary Track akses materi peserta
// @Description Mencatat bahwa peserta telah membuka materi
// @Tags Materi Progress
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id_materi path int true "ID Materi"
// @Router /materi/peserta/{id_materi}/progress [post]
func (h *Handler) OpenMateri(c *fiber.Ctx) error {

	userID, ok := c.Locals("sub").(int)
	if !ok {
		return response.Error(c, 401, "unauthorized", "UNAUTHORIZED", nil)
	}

	materiID, err := c.ParamsInt("id_materi")
	if err != nil || materiID <= 0 {
		return response.Error(c, 400, "invalid materi id", "BAD_REQUEST", nil)
	}

	if err := h.Service.OpenMateri(c.Context(), userID, materiID); err != nil {
		return response.Error(c, 500, err.Error(), "INTERNAL_ERROR", nil)
	}

	return response.Success(c, nil)
}


// GetMateriProgress godoc
// @Summary Get progress materi peserta
// @Description Mengambil raw progress materi milik peserta yang sedang login
// @Tags Materi Progress
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Nomor halaman" default(1)
// @Param limit query int false "Jumlah data per halaman" default(20)
// @Param id_modul query int false "Filter berdasarkan ID Modul"
// @Router /materi/peserta/progress [get]
func (h *Handler) GetMateriProgress(c *fiber.Ctx) error {

	userID, ok := c.Locals("sub").(int)
	if !ok {
		return response.Error(c, 401, "unauthorized", "UNAUTHORIZED", nil)
	}

	page, limit := c.QueryInt("page", 1), c.QueryInt("limit", 20)
	if page < 1 {
		return response.Error(c, 400, "page must be greater than 0", "BAD_REQUEST", nil)
	}
	if limit < 1 {
		return response.Error(c, 400, "limit must be greater than 0", "BAD_REQUEST", nil)
	}
	if limit > 100 {
		limit = 100
	}

	var modulID *int
	if idModul := c.Query("id_modul"); idModul != "" {
		id, err := strconv.Atoi(idModul)
		if err != nil || id <= 0 {
			return response.Error(c, 400, "invalid modul id", "BAD_REQUEST", nil)
		}
		modulID = &id
	}

	data, total, err := h.Service.GetMateriProgress(c.Context(), userID, modulID, page, limit)
	if err != nil {
		return response.Error(c, 500, err.Error(), "INTERNAL_ERROR", nil)
	}

	totalPages := 0
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(limit)))
	}

	result := MateriProgressListResponseDTO{
		Data: data,
		Pagination: PaginationDTO{
			Page: page, Limit: limit, Total: total, TotalPages: totalPages,
		},
	}

	return response.Success(c, result)
}


// GetMateriProgressMonitoring godoc
// @Summary Get monitoring progress materi peserta
// @Description Mengambil ringkasan progress materi peserta yang sedang login
// @Tags Materi Progress
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id_modul query int false "Filter berdasarkan ID Modul"
// @Router /materi/peserta/progress/monitoring [get]
func (h *Handler) GetMateriProgressMonitoring(c *fiber.Ctx) error {

	userID, ok := c.Locals("sub").(int)
	if !ok {
		return response.Error(c, 401, "unauthorized", "UNAUTHORIZED", nil)
	}

	var modulID *int
	if idModul := c.Query("id_modul"); idModul != "" {
		id, err := strconv.Atoi(idModul)
		if err != nil || id <= 0 {
			return response.Error(c, 400, "invalid modul id", "BAD_REQUEST", nil)
		}
		modulID = &id
	}

	data, err := h.Service.GetMateriProgressMonitoring(c.Context(), userID, modulID)
	if err != nil {
		return response.Error(c, 500, err.Error(), "INTERNAL_ERROR", nil)
	}

	return response.Success(c, data)
}