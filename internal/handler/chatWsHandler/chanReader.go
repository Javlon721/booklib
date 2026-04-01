package chatWsHandler

import (
	"context"
	"encoding/json"
	"io"
)

type chanReader struct {
	ctx context.Context
	pw  *io.PipeWriter
	pr  *io.PipeReader
}

func NewChanReader(ctx context.Context) *chanReader {
	pr, pw := io.Pipe()
	return &chanReader{
		ctx: ctx,
		pr:  pr,
		pw:  pw,
	}
}

// make sure that rch should return valid json encodable data
func (c *chanReader) Listen(rch <-chan any) {
	defer c.pw.Close()

	encoder := json.NewEncoder(c.pw)

	for {
		select {
		case <-c.ctx.Done():
			return
		case value, ok := <-rch:
			if !ok {
				return
			}

			_ = encoder.Encode(value)
		}
	}
}

func (c *chanReader) Read(b []byte) (n int, err error) {
	return c.pr.Read(b)
}
