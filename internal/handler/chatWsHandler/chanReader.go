package chatWsHandler

import (
	"context"
	"io"
)

type chanReader struct {
	ctx context.Context
	rch <-chan string
}

func NewChanReader(ctx context.Context, rch <-chan string) *chanReader {
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

		n = copy(p, data)
	}

	return
}
