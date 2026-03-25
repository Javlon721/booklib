package topichandler

import (
	topicbus "github.com/Javlon721/booklib/internal/bus/topicBus"
	"github.com/Javlon721/booklib/internal/handler/errs"
	"github.com/Javlon721/booklib/internal/types/name"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type queryParams struct {
	CreatedBy string `query:"created_by"`
	Title     string `query:"title"`
}

func parseQueryParams(c fiber.Ctx) (queryParams, error) {
	var queries queryParams

	if err := c.Bind().Query(&queries); err != nil {
		return queryParams{}, err
	}

	return queries, nil
}

func parseFilter(queries queryParams) (topicbus.QueryFilter, error) {
	var fieldErrors errs.FieldErrors
	var filter topicbus.QueryFilter

	if queries.Title != "" {
		title, err := name.Parse(queries.Title)
		switch err {
		case nil:
			filter.Title = &title
		default:
			fieldErrors.Add("title", err)
		}
	}

	if queries.CreatedBy != "" {
		createdBy, err := uuid.Parse(queries.CreatedBy)
		switch err {
		case nil:
			filter.CreatedBy = &createdBy
		default:
			fieldErrors.Add("createdBy", err)
		}
	}

	if len(fieldErrors) > 0 {
		return topicbus.QueryFilter{}, fieldErrors.ToError()
	}

	return filter, nil
}
