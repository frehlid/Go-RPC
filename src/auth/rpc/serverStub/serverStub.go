package serverStub

import (
	"encoding/json"

	"local/auth"
	"local/auth/rpc/api"
	"local/lib/helpers"
)

func Signup(args []byte) []byte {
	var signupArgs api.SignupArgs
	err := json.Unmarshal(args, &signupArgs)
	helpers.CheckForError(err)

	result := auth.Signup(
		signupArgs.Id,
		signupArgs.Password,
	)

	resultStruct := api.SignupResult{Result: result}
	resultBuf, err := json.Marshal(resultStruct)
	helpers.CheckForError(err)
	return resultBuf
}

func Login(data []byte) []byte {
	var args api.LoginArgs
	err := json.Unmarshal(data, &args)
	helpers.CheckForError(err)

	result := auth.Login(
		args.Id,
		args.Password,
	)

	resultStruct := api.LoginResult{Cap: result}
	resultBuf, err := json.Marshal(resultStruct)
	helpers.CheckForError(err)
	return resultBuf
}

func GetId(data []byte) []byte {
	var args api.GetIdArgs
	helpers.VerbosePrint("Args to getID:", string(data))

	err := json.Unmarshal(data, &args)
	helpers.CheckForError(err)

	result := auth.GetId(
		args.Cap,
	)

	helpers.VerbosePrint("GET ID AUTH SERVER STUB :", result)

	resultStruct := api.GetIdResult{result}
	resultBuf, err := json.Marshal(resultStruct)
	helpers.CheckForError(err)

	return resultBuf
}
