package user

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(r fiber.Router, handler *Handler) {
	user := r.Group("/user")
	
	user.Get("/me", handler.Me)
	user.Put("/change-password", handler.ChangePassword)
}