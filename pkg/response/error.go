package response

import (
	"github.com/gofiber/fiber/v2"

	"root/pkg/config"
)

func Error(ctx *fiber.Ctx, status int, err error, message string) error {
	cfg := config.GetConfig()
	errorRes := map[string]interface{}{
		"message": message,
	}

	if cfg.Enviroment != config.ProductionEnv {
		errorRes["debug"] = err.Error()
	}

	return ctx.Status(status).JSON(Response{
		Error: errorRes,
	})
}


type ErrorResponse struct {
    StatusCode int
    Message    string
    Err        error
}

func (e *ErrorResponse) Error() string {
    return e.Message
}


func GetErroField(err error) (int, string) {
	if errWrap, ok := err.(*ErrorResponse); ok {
		return errWrap.StatusCode, errWrap.Message
	}
	return 500, "Something went wrong"
}
