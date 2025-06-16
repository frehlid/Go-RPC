package main

import (
	"context"
	auth "local/auth/rpc/clientStub"
	db "local/db/rpc/clientStub"
	"local/lib/finalizer"
	"local/lib/handler"
	"local/lib/helpers"
	"local/lib/transport"
	messageServerStub "local/message/rpc/serverStub"
	"os"
)

var messageFuncMap = map[string]func(args []byte) []byte{
	"Send":              messageServerStub.Send,
	"Receive":           messageServerStub.Receive,
	"SetSendingAllowed": messageServerStub.SetSendingAllowed,
	"SetReceiver":       messageServerStub.SetReceiver,
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
	// bind db to addr given in args
	args := os.Args[1:]
	dbAddr := args[0]
	db.Bind(transport.NewNetAddr(dbAddr))

	// bind auth to the addr stored in nameserver, loop until db gives response
	var authAddr = ""
	for {
		authAddr = db.Get("auth")
		helpers.VerbosePrint("MESSAGED got auth address:", authAddr)
		if authAddr != "" {
			break
		}
	}
	auth.Bind(transport.NewNetAddr(authAddr))
	go transport.Listen(ctx, messageFuncMap, handler.Handler)

	listenAddr := transport.LocalAddr()
	db.Put("message", listenAddr)
}
