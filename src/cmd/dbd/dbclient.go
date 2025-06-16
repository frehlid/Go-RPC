package main

import (
	"local/db/rpc/clientStub"
	"local/lib/transport"
)

func main() {
	addr := transport.NewNetAddr("127.0.0.1:" + "58134")
	clientStub.Bind(addr)

	putResult := clientStub.Put("Who is a poopy face?", "Dieter + Matias")
	helpers.VerbosePrint(putResult)
	getResult := clientStub.Get("Who is a poopy face?")
	helpers.VerbosePrint(getResult)
}
