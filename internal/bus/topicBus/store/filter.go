package store

import (
	"bytes"
	"fmt"
	"strings"

	topicbus "github.com/Javlon721/booklib/internal/bus/topicBus"
	"github.com/jackc/pgx/v5"
)

func applyFilter(filter topicbus.QueryFilter, buf *bytes.Buffer) pgx.NamedArgs {
	var wc []string
	namedArgs := make(pgx.NamedArgs)

	if filter.CreatedBy != nil {
		wc = append(wc, "created_by = @created_by")
		namedArgs["created_by"] = filter.CreatedBy
	}

	if filter.Title != nil {
		namedArgs["title"] = fmt.Sprintf("%%%s%%", filter.Title)
		wc = append(wc, "title LIKE @title")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}

	return namedArgs
}
