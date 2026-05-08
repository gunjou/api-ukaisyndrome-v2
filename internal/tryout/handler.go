package tryout

import (
	"api-ukaisyndrome-v2/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Service *Service
}

// GetTryoutPeserta godoc
// @Summary Get tryout peserta
// @Description Get list tryout berdasarkan kelas peserta
// @Tags Tryout
// @Accept json
// @Produce json
// @Security BearerAuth
// @Router /tryout/peserta [get]
func (h *Handler) GetTryoutPeserta(c *fiber.Ctx) error {

	userID, ok := c.Locals("sub").(int)
	if !ok {
		return response.Error(c, 401, "unauthorized", "UNAUTHORIZED", nil)
	}

	data, err := h.Service.GetTryoutPeserta(c.Context(), userID)
	if err != nil {
		return response.Error(c, 500, err.Error(), "INTERNAL_ERROR", nil)
	}

	return response.Success(c, data)
}


// StartTryout godoc
// @Summary Start tryout
// @Description Memulai tryout peserta
// @Tags Tryout
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id_tryout path int true "ID Tryout"
// @Router /tryout/{id_tryout}/start [post]
func (h *Handler) StartTryout(c *fiber.Ctx) error {

	userID := c.Locals("sub").(int)

	tryoutID, err := c.ParamsInt("id_tryout")
	if err != nil {
		return response.Error(c, 400, "invalid tryout id", "BAD_REQUEST", nil)
	}

	res, err := h.Service.StartTryout(c.Context(), userID, tryoutID)
	if err != nil {
		return response.Error(c, 400, err.Error(), "BAD_REQUEST", nil)
	}

	return response.Success(c, res)
}


// GetSoalTryout godoc
// @Summary Get soal tryout
// @Description Ambil semua soal untuk attempt
// @Tags Tryout
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param attempt_token path string true "Attempt Token"
// @Router /tryout/attempt/{attempt_token} [get]
func (h *Handler) GetSoalTryout(c *fiber.Ctx) error {

	userID := c.Locals("sub").(int)

	token := c.Params("attempt_token")

	data, err := h.Service.GetSoalTryout(
		c.Context(),
		userID,
		token,
	)

	if err != nil {
		return response.Error(c, 400, err.Error(), "BAD_REQUEST", nil)
	}

	return response.Success(c, data)
}



// SaveAnswers godoc
// @Summary Autosave jawaban tryout
// @Description Menyimpan jawaban peserta secara berkala
// @Tags Tryout
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param attempt_token path string true "Attempt Token"
// @Param request body SaveAnswerRequest true "Answers"
// @Router /tryout/attempt/{attempt_token}/answers [put]
func (h *Handler) SaveAnswers(c *fiber.Ctx) error {

	userID := c.Locals("sub").(int)
	token := c.Params("attempt_token")

	var req SaveAnswerRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, 400, "invalid request body", "BAD_REQUEST", nil)
	}

	if len(req.Answers) == 0 {
		return response.Error(c, 400, "answers cannot be empty", "BAD_REQUEST", nil)
	}

	err := h.Service.SaveAnswers(
		c.Context(),
		userID,
		token,
		req.Answers,
	)

	if err != nil {
		return response.Error(c, 400, err.Error(), "BAD_REQUEST", nil)
	}

	return response.Success(c, "saved")
}


// SubmitTryout godoc
// @Summary Submit tryout
// @Description Submit hasil tryout peserta
// @Tags Tryout
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param attempt_token path string true "Attempt Token"
// @Router /tryout/attempt/{attempt_token}/submit [post]
func (h *Handler) SubmitTryout(c *fiber.Ctx) error {

	userID := c.Locals("sub").(int)
	token := c.Params("attempt_token")

	data, err := h.Service.SubmitTryout(c.Context(), userID, token)
	if err != nil {
		return response.Error(c, 400, err.Error(), "BAD_REQUEST", nil)
	}

	return response.Success(c, data)
}


// ResumeTryout godoc
// @Summary Resume tryout
// @Description Melanjutkan attempt yang sedang berjalan
// @Tags Tryout
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param attempt_token path string true "Attempt Token"
// @Router /tryout/attempt/{attempt_token}/resume [get]
func (h *Handler) ResumeTryout(c *fiber.Ctx) error {

	userID := c.Locals("sub").(int)
	token := c.Params("attempt_token")

	data, err := h.Service.ResumeAttempt(
		c.Context(),
		userID,
		token,
	)

	if err != nil {
		return response.Error(c, 400, err.Error(), "BAD_REQUEST", nil)
	}

	return response.Success(c, data)
}


// GetReports godoc
// @Summary Get laporan tryout
// @Tags Tryout
// @Security BearerAuth
// @Router /tryout/report [get]
func (h *Handler) GetReports(c *fiber.Ctx) error {

	userID := c.Locals("sub").(int)

	data, err := h.Service.GetReports(c.Context(), userID)
	if err != nil {
		return response.Error(c, 500, err.Error(), "INTERNAL_ERROR", nil)
	}

	return response.Success(c, data)
}


// GetReview godoc
// @Summary Get pembahasan tryout
// @Tags Tryout
// @Security BearerAuth
// @Param attempt_token path string true "Attempt Token"
// @Router /tryout/report/{attempt_token} [get]
func (h *Handler) GetReview(c *fiber.Ctx) error {

	userID := c.Locals("sub").(int)
	token := c.Params("attempt_token")

	data, err := h.Service.GetReview(c.Context(), userID, token)
	if err != nil {
		return response.Error(c, 500, err.Error(), "INTERNAL_ERROR", nil)
	}

	return response.Success(c, data)
}