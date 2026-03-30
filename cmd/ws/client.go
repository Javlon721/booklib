package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/coder/websocket"
)

var (
	userToken = flag.String("token", "", "user token for authorization")
	url       = flag.String("url", "", "url for ws connection")
)

func main() {
	slog.Info("startup", "status", "initizlizing config")

	flag.Parse()

	if *userToken == "" {
		panic("token for authorization not provided")
	}

	if *url == "" {
		panic("url for ws connection not provided")
	}

	headers := http.Header{}

	headers.Add("Authorization", fmt.Sprintf("Bearer %s", *userToken))

	wsOptions := &websocket.DialOptions{HTTPHeader: headers}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()

	// -------------------------------------------------------------------------

	slog.Info("startup", "status", "establishing connection")

	conn, resp, err := websocket.Dial(ctx, *url, wsOptions)

	if err != nil {
		body := map[string]string{}

		if resp != nil {
			_ = json.NewDecoder(resp.Body).Decode(&body)
		}

		slog.Error("while establishing conn", "err", err, slog.Any("resp", body))

		return
	}

	defer conn.CloseNow()

	// -------------------------------------------------------------------------

	readerFn := func(ctx context.Context, conn *websocket.Conn) <-chan string {
		data := make(chan string)

		go func() {
			defer close(data)

			mt, msg, err := conn.Read(ctx)

			if err != nil {
				slog.Error("while reading", "err", err)
				return
			}

			if mt != websocket.MessageText {
				slog.Error("invalid message type", "type", mt)
				return
			}

			data <- string(msg)
		}()

		return data
	}

	slog.Info("startup", "status", "get read chanel")

	data := readerFn(context.Background(), conn)

	time.Sleep(time.Second)

	err = conn.Write(context.Background(), websocket.MessageText, []byte("works"))

	if err != nil {
		slog.Error("while writing message", "err", err)
		return
	}

	fmt.Printf("Data recieved: %s\n", <-data)

	conn.Close(websocket.StatusNormalClosure, "")
}
