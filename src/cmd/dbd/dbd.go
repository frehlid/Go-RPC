package main

import (
	"context"
	"fmt"
	"local/db/rpc/serverStub"
	"local/lib/finalizer"
	"local/lib/handler"
	"local/lib/transport"
)

var dbFuncMap = map[string]func(args []byte) []byte{
	"Put": serverStub.Put,
	"Get": serverStub.Get,
}

func main() {
	ctx, cancel := finalizer.WithCancel(context.Background())
	defer cancel()

	startListening(ctx)

	select {
	case <-ctx.Done(): // ctx cancelled
		return
	}
}

func startListening(ctx context.Context) {
	go transport.Listen(ctx, dbFuncMap, handler.Handler)

	listenAddr := transport.LocalAddr()
	fmt.Println("Listening on " + listenAddr)
}
