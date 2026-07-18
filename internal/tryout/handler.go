package tryout

import (
	"api-ukaisyndrome-v2/pkg/response"
	"log"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Service *Service
}



/* ========================================================================== */
/*                  //SECTION - ENDPOINT LIST TRYOUT SECTION                  */
/* ========================================================================== */

//ANCHOR - GET TRYOUT PESERTA
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


//ANCHOR - START TRYOUT
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


//ANCHOR - GET SOAL TRYOUT
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


//ANCHOR - SAVE ANSWERS
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

	// if len(req.Answers) == 0 {
	// 	return response.Error(c, 400, "answers cannot be empty", "BAD_REQUEST", nil)
	// }

	err := h.Service.SaveAnswers(
		c.Context(),
		userID,
		token,
		req.Answers,
	)

	if err != nil {
		log.Println("SAVE ANSWER ERROR:", err)
		return response.Error(c, 400, err.Error(), "BAD_REQUEST", nil)
	}

	return response.Success(c, "saved")
}


//ANCHOR - SUBMIT TRYOUT
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


//ANCHOR - RESUME TRYOUT
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

//ANCHOR - GET ONGOING TRYOUT
// GetOngoingTryout godoc
// @Summary Get ongoing tryout
// @Description Mengecek apakah peserta memiliki tryout yang masih ongoing atau sudah expired
// @Tags Tryout
// @Accept json
// @Produce json
// @Security BearerAuth
// @Router /tryout/is_ongoing [get]
func (h *Handler) GetOngoingTryout(c *fiber.Ctx) error {

	userID := c.Locals("sub").(int)

	data, err := h.Service.GetOngoingTryout(
		c.Context(),
		userID,
	)
	if err != nil {
		return response.Error(
			c,
			500,
			err.Error(),
			"INTERNAL_ERROR",
			nil,
		)
	}

	return response.Success(
		c,
		data,
	)
}
/* =========================== //!SECTION - TRYOUT ========================== */


/* ========================================================================== */
/*                  //SECTION - ENDPOINT LIST REPORT SECTION                  */
/* ========================================================================== */

//ANCHOR - GET REPORTS
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

//ANCHOR - GET REVIEW
// GetReview godoc
// @Summary Get pembahasan tryout
// @Tags Tryout
// @Security BearerAuth
// @Param attempt_token path string true "Attempt Token"
// @Router /tryout/report/{attempt_token} [get]
func (h *Handler) GetReview(c *fiber.Ctx) error {

	userID := c.Locals("sub").(int)
	token := c.Params("attempt_token")

	title, data, err := h.Service.GetReview(
		c.Context(),
		userID,
		token,
	)

	if err != nil {
		return response.Error(
			c,
			500,
			err.Error(),
			"INTERNAL_ERROR",
			nil,
		)
	}

	return response.SuccessWithAdditional(
		c,
		data,
		map[string]interface{}{
			"title": title,
		},
	)
}
/* =========================== //!SECTION - REPORT ========================== */


/* ========================================================================== */
/*                //SECTION - ENDPOINT LIST LEADERBOARD SECTION               */
/* ========================================================================== */

//ANCHOR - GLOBAL LEADERBOARD
// GetGlobalLeaderboard godoc
// @Summary Get global leaderboard
// @Description Leaderboard seluruh peserta berdasarkan tryout
// @Tags Tryout
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id_tryout path int true "ID Tryout"
// @Router /tryout/{id_tryout}/leaderboard [get]
func (h *Handler) GetGlobalLeaderboard(c *fiber.Ctx) error {

	userID := c.Locals("sub").(int)

	tryoutID, err := c.ParamsInt("id_tryout")
	if err != nil {
		return response.Error(
			c,
			400,
			"invalid tryout id",
			"BAD_REQUEST",
			nil,
		)
	}

	data, err := h.Service.GetGlobalLeaderboard(
		c.Context(),
		userID,
		tryoutID,
	)

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


//ANCHOR - CLASS LEADERBOARD
// GetClassLeaderboard godoc
// @Summary Get class leaderboard
// @Description Leaderboard berdasarkan kelas
// @Tags Tryout
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id_tryout path int true "ID Tryout"
// @Param class_id query int false "Class ID (Admin & Mentor only)"
// @Router /tryout/{id_tryout}/leaderboard/class [get]
func (h *Handler) GetClassLeaderboard(c *fiber.Ctx) error {

	userID := c.Locals("sub").(int)
	role := c.Locals("role").(string)

	tryoutID, err := c.ParamsInt("id_tryout")
	if err != nil {
		return response.Error(
			c,
			400,
			"invalid tryout id",
			"BAD_REQUEST",
			nil,
		)
	}

	// optional (hanya digunakan admin & mentor)
	classID := c.QueryInt("class_id")

	data, err := h.Service.GetClassLeaderboard(
		c.Context(),
		userID,
		role,
		tryoutID,
		classID,
	)

	if err != nil {
		return response.Error(
			c,
			400,
			err.Error(),
			"BAD_REQUEST",
			nil,
		)
	}

	return response.Success(c, data)
}
/* ======================== //!SECTION - LEADERBOARD ======================== */



/* ========================================================================== */
/*                 //SECTION - ENDPOINT LIST ANALYTICS SECTION                */
/* ========================================================================== */

//ANCHOR - ANALYTICS
// GetStatistics godoc
// @Summary Get tryout statistics
// @Description Statistik tryout (global / class)
// @Tags Tryout
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id_tryout path int true "ID Tryout"
// @Param class_id query int false "Class ID (Admin & Mentor only)"
// @Router /tryout/{id_tryout}/statistics [get]
func (h *Handler) GetStatistics(c *fiber.Ctx) error {

	userID := c.Locals("sub").(int)
	role := c.Locals("role").(string)

	tryoutID, err := c.ParamsInt("id_tryout")
	if err != nil {
		return response.Error(
			c,
			400,
			"invalid tryout id",
			"BAD_REQUEST",
			nil,
		)
	}

	classID := c.QueryInt("class_id")

	data, err := h.Service.GetStatistics(
		c.Context(),
		userID,
		role,
		tryoutID,
		classID,
	)
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

//ANCHOR - GET DISTRIBUTION
// GetDistribution godoc
// @Summary Get score distribution
// @Description Distribusi nilai tryout
// @Tags Tryout
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id_tryout path int true "ID Tryout"
// @Param class_id query int false "Class ID"
// @Router /tryout/{id_tryout}/distribution [get]
func (h *Handler) GetDistribution(c *fiber.Ctx) error {

	userID := c.Locals("sub").(int)
	role := c.Locals("role").(string)

	tryoutID, err := c.ParamsInt("id_tryout")
	if err != nil {
		return response.Error(c, 400, "invalid tryout id", "BAD_REQUEST", nil)
	}

	classID := c.QueryInt("class_id")

	data, err := h.Service.GetDistribution(
		c.Context(),
		userID,
		role,
		tryoutID,
		classID,
	)

	if err != nil {
		return response.Error(c, 500, err.Error(), "INTERNAL_ERROR", nil)
	}

	return response.Success(c, data)
}


//ANCHOR - GET QUESTION ANALYSIS
// GetQuestionAnalysis godoc
// @Summary Get question analysis
// @Description Get analysis for each question in a tryout
// @Tags Tryout
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id_tryout path int true "Tryout ID"
// @Param class_id query int false "Class ID (optional for admin & mentor)"
// @Param sort query string false "Sort: default | hardest | easiest"
// @Router /tryout/{id_tryout}/questions [get]
func (h *Handler) GetQuestionAnalysis(c *fiber.Ctx) error {

	userID := c.Locals("sub").(int)
	role := c.Locals("role").(string)

	tryoutID, err := c.ParamsInt("id_tryout")
	if err != nil {
		return response.Error(
			c,
			400,
			"invalid tryout id",
			"BAD_REQUEST",
			nil,
		)
	}

	// optional
	classID := c.QueryInt("class_id")

	// optional
	sortBy := c.Query("sort", "default")

	data, err := h.Service.GetQuestionAnalysis(
		c.Context(),
		userID,
		role,
		tryoutID,
		classID,
		sortBy,
	)
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


//ANCHOR - GET QUESTION CHOICE ANALYSIS
// GetQuestionChoices godoc
// @Summary Get question choice analysis
// @Description Get answer distribution for a question
// @Tags Tryout
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id_soal path int true "Question ID"
// @Param class_id query int false "Class ID (Optional for admin & mentor)"
// @Router /tryout/question/{id_soal}/choices [get]
func (h *Handler) GetQuestionChoices(c *fiber.Ctx) error {

	userID := c.Locals("sub").(int)
	role := c.Locals("role").(string)

	questionID, err := c.ParamsInt("id_soal")
	if err != nil {
		return response.Error(
			c,
			400,
			"invalid question id",
			"BAD_REQUEST",
			nil,
		)
	}

	// optional
	classID := c.QueryInt("class_id")

	data, err := h.Service.GetQuestionChoices(
		c.Context(),
		userID,
		role,
		questionID,
		classID,
	)
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
/* ========================= //!SECTION - ANALYTICS ========================= */
