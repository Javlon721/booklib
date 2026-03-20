package middleware

import (
	"errors"
	"log/slog"
	"path"

	"github.com/Javlon721/booklib/internal/handler/errs"
	"github.com/gofiber/fiber/v3"
)

func Errors(logger *slog.Logger) fiber.ErrorHandler {
	return func(c fiber.Ctx, err error) error {
		var e *errs.Error

		if !errors.As(err, &e) {
			e = errs.New(errs.Internal, err)
		}

		logger.Error("handled error during request",
			"err", e,
			"source_err_file", path.Base(e.FileName),
			"source_err_func", path.Base(e.FuncName),
		)

		if e.Code.Equal(errs.Internal) {
			return c.Status(e.HTTPStatus()).JSON(map[string]string{"error": "internal server error"})
		}

		return c.Status(e.HTTPStatus()).JSON(map[string]string{"error": e.Message})
	}
}
