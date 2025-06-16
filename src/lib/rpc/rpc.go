package rpc

import (
	"bytes"
	"encoding/json"
	"local/lib/helpers"
	"local/lib/transport"
)

type RPCData struct {
	Method string
	Data   []byte
}

func MarshalArgs(args RPCData) *bytes.Buffer {
	buff, err := json.Marshal(args)
	helpers.CheckForError(err)
	return bytes.NewBuffer(buff)
}

func UnmarshalArgs(buffer *bytes.Buffer) RPCData {
	var args RPCData
	//helpers.VerbosePrint(buffer.String())
	err := json.Unmarshal(buffer.Bytes(), &args)
	helpers.CheckForError(err)
	return args
}

func RPCCall(to transport.NetAddr, method string, args []byte) []byte {
	rpcArgs := RPCData{method, args}
	buf := MarshalArgs(rpcArgs)

	rBuf, err := transport.Call(buf, to)
	helpers.CheckForError(err)

	return rBuf.Bytes()
}
