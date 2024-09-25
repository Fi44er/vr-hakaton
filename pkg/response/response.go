package response

import (
	"github.com/gofiber/fiber/v2"
)

// Response schema
type Response struct {
	Result interface{} `json:"result"`
	Error  interface{} `json:"error"`
}

func JSON(ctx *fiber.Ctx, status int, data interface{}) error {
  return ctx.Status(status).JSON(Response{
    Result: data,
  })
}
