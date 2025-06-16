package main

import (
	"context"
	authServerStub "local/auth/rpc/serverStub"
	dbClientStub "local/db/rpc/clientStub"
	"local/lib/finalizer"
	handler "local/lib/handler"
	"local/lib/transport"
	"os"
)

var authFuncMap = map[string]func(args []byte) []byte{
	"Signup": authServerStub.Signup,
	"Login":  authServerStub.Login,
	"GetId":  authServerStub.GetId,
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
	args := os.Args[1:]
	dbAddr := args[0]

	// bind db to the addr given from args
	dbClientStub.Bind(transport.NewNetAddr(dbAddr))

	// start listening
	go transport.Listen(ctx, authFuncMap, handler.Handler)

	// once listen has an addr, send it to the nameserver
	localAddr := transport.LocalAddr()
	dbClientStub.Put("auth", localAddr)
}
