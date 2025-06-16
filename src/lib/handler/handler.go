package handler

import (
	"bytes"
	"encoding/json"
	"local/lib/helpers"
	"local/lib/rpc"
	"local/lib/transport"
)

func Handler(msg *bytes.Buffer, from transport.NetAddr, funcMap transport.FunctionMap) {
	args := rpc.UnmarshalArgs(msg)

	funcName := args.Method
	fn, ok := funcMap[funcName]

	var result []byte
	var err error
	if ok {
		result = fn(args.Data)
	} else {
		result, err = json.Marshal("Error executing " + funcName)
		helpers.CheckForError(err)
	}

	transport.Reply(bytes.NewBuffer(result), from)
}
