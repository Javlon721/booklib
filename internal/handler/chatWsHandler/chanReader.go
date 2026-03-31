package chatWsHandler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
)

type chanReader struct {
	ctx context.Context
	rch <-chan any
}

func NewChanReader(ctx context.Context, rch <-chan any) *chanReader {
	return &chanReader{
		ctx: ctx,
		rch: rch,
	}
}

func (c *chanReader) Read(p []byte) (n int, err error) {
	select {
	case <-c.ctx.Done():
		return 0, io.EOF

	case data, ok := <-c.rch:
		if !ok {
			return 0, io.EOF
		}

		b, err := json.Marshal(data)

		if err != nil {
			return 0, err
		}

		n = copy(p, fmt.Append(b, "\n"))
	}

	return
}
