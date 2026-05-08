package tryout

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(r fiber.Router, handler *Handler) {

	tryout := r.Group("/tryout")

	tryout.Get("/peserta", handler.GetTryoutPeserta)
	tryout.Post("/:id_tryout/start", handler.StartTryout)
	tryout.Get("/attempt/:attempt_token", handler.GetSoalTryout)
	tryout.Put("/attempt/:attempt_token/answers", handler.SaveAnswers)
	tryout.Post("/attempt/:attempt_token/submit", handler.SubmitTryout)
	tryout.Get("/attempt/:attempt_token/resume", handler.ResumeTryout)
	tryout.Get("/report", handler.GetReports)
	tryout.Get("/report/:attempt_token", handler.GetReview)
}