package response

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type Meta struct {
	ResponseTimeMs int64  `json:"response_time_us"`
	Timestamp      string `json:"timestamp"`
	RequestID      string `json:"request_id,omitempty"`
}

type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Meta    Meta        `json:"meta"`
}

// ==========================
// HELPER FUNCTIONS
// ==========================
func getResponseTime(c *fiber.Ctx) int64 {

	start, ok := c.Locals("start_time").(time.Time)
	if !ok {
		return 0
	}

	return time.Since(start).Milliseconds()
}

func buildMeta(c *fiber.Ctx, rt int64) Meta {

	reqID, _ := c.Locals("request_id").(string)

	return Meta{
		ResponseTimeMs: rt,
		Timestamp:      time.Now().UTC().Format(time.RFC3339),
		RequestID:      reqID,
	}
}

// ==========================
// SUCCESS
// ==========================
func Success(c *fiber.Ctx, data interface{}) error {

	rt := getResponseTime(c)

	return c.JSON(APIResponse{
		Status:  "success",
		Message: "OK",
		Data:    data,
		Meta:    buildMeta(c, rt),
	})
}

// ==========================
// ERROR
// ==========================
func Error(c *fiber.Ctx, statusCode int, message string, code string, details interface{}) error {

	rt := getResponseTime(c)

	return c.Status(statusCode).JSON(APIResponse{
		Status:  "error",
		Message: message,
		Error: map[string]interface{}{
			"code":    code,
			"details": details,
		},
		Meta: buildMeta(c, rt),
	})
}