package clientStub

import (
	"encoding/json"

	"local/db/rpc/api"
	"local/lib/helpers"
	"local/lib/rpc"
	"local/lib/transport"
)

func Bind(toAddr transport.NetAddr) {
	to = toAddr
}

func Put(key string, value string) bool {
	args, err := json.Marshal(api.PutArgs{
		Key: key, Value: value,
	})
	helpers.CheckForError(err)

	result := rpc.RPCCall(
		to,
		"Put",
		args,
	)

	helpers.VerbosePrint(result)

	var putRes api.PutResult
	err = json.Unmarshal(result, &putRes)
	helpers.CheckForError(err)

	return putRes.Result
}

func Get(key string) string {
	args, err := json.Marshal(api.GetArgs{
		Key: key,
	})
	helpers.CheckForError(err)
	result := rpc.RPCCall(
		to,
		"Get",
		args,
	)

	var getRes api.GetResult
	err = json.Unmarshal(result, &getRes)
	helpers.CheckForError(err)

	return getRes.Result
}

var to transport.NetAddr
