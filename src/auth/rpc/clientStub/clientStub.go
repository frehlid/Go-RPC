package clientStub

import (
	"encoding/json"

	"local/auth/rpc/api"
	"local/auth/types"
	"local/lib/helpers"
	"local/lib/rpc"
	"local/lib/transport"
)

func Bind(toAddr transport.NetAddr) {
	to = toAddr
}

func Signup(id string, password string) bool {
	argStruct := api.SignupArgs{id, password}
	args, err := json.Marshal(argStruct)
	helpers.CheckForError(err)

	result := rpc.RPCCall(
		to,
		"Signup",
		args,
	)

	var signupResult api.SignupResult
	err = json.Unmarshal(result, &signupResult)
	helpers.CheckForError(err)

	return signupResult.Result
}

func Login(id string, password string) types.UserCap {
	args, err := json.Marshal(
		api.LoginArgs{
			Id: id, Password: password,
		})
	helpers.CheckForError(err)
	result := rpc.RPCCall(to,
		"Login",
		args,
	)
	var loginResult api.LoginResult
	err = json.Unmarshal(result, &loginResult)
	helpers.CheckForError(err)

	return loginResult.Cap
}

func GetId(cap types.UserCap) string {
	args, err := json.Marshal(
		api.GetIdArgs{
			Cap: cap,
		})
	helpers.CheckForError(err)

	result := rpc.RPCCall(
		to,
		"GetId",
		args,
	)

	helpers.VerbosePrint("GET ID AUTH CLIENT STUB: ", result)

	var getIdResult api.GetIdResult
	err = json.Unmarshal(result, &getIdResult)
	helpers.CheckForError(err)

	return getIdResult.Id
}

var to transport.NetAddr
