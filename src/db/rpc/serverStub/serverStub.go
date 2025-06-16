package serverStub

import (
	"encoding/json"
	"local/db"
	"local/db/rpc/api"
	"local/lib/helpers"
)

func Put(data []byte) []byte {
	var args api.PutArgs
	err := json.Unmarshal(data, &args)
	helpers.CheckForError(err)

	db.Put(
		args.Key,
		args.Value,
	)

	resultBuf := api.PutResult{Result: true}
	result, err := json.Marshal(resultBuf)
	helpers.CheckForError(err)
	return result
}

func Get(data []byte) []byte {
	var args api.GetArgs
	err := json.Unmarshal(data, &args)
	helpers.CheckForError(err)

	getResult := db.Get(args.Key)

	result := api.GetResult{Result: getResult}
	ret, err := json.Marshal(result)
	helpers.CheckForError(err)

	return ret
}
