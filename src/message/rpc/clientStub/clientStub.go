package clientStub

import (
	"encoding/json"
	authTypes "local/auth/types"
	"local/lib/helpers"
	"local/lib/rpc"
	"local/lib/transport"
	"local/message/rpc/api"
	"local/message/rpc/proxy"
	"local/message/types"
)

func Bind(to transport.NetAddr) {
	toAddr = to
}

func Send(user authTypes.UserCap, to string, text string) bool {
	argStruct := api.SendArgs{user, to, text}
	args, err := json.Marshal(argStruct)
	helpers.CheckForError(err)

	result := rpc.RPCCall(
		toAddr,
		"Send",
		args,
	)

	var sendResult api.SendResult
	err = json.Unmarshal(result, &sendResult)
	helpers.CheckForError(err)

	return sendResult.Result
}

func Receive(user authTypes.UserCap) *types.Message {
	argStruct := api.ReceiveArgs{user}
	args, err := json.Marshal(argStruct)
	helpers.CheckForError(err)

	result := rpc.RPCCall(
		toAddr,
		"Receive",
		args,
	)

	var receiveResult api.ReceiveResult
	err = json.Unmarshal(result, &receiveResult)
	helpers.CheckForError(err)

	return receiveResult.Message
}

func SetSendingAllowed(user authTypes.UserCap, from string, allowed bool) bool {
	argStruct := api.SetSendingAllowedArgs{user, from, allowed}
	args, err := json.Marshal(argStruct)
	helpers.CheckForError(err)

	result := rpc.RPCCall(
		toAddr,
		"SetSendingAllowed",
		args,
	)

	var setSendingAllowedResult api.SetSendingAllowedResult
	err = json.Unmarshal(result, &setSendingAllowedResult)
	helpers.CheckForError(err)

	return setSendingAllowedResult.Result
}

func SetReceiver(user authTypes.UserCap, receiver types.Receiver, receive bool) bool {
	receiverId := proxy.LocalObjectToRemoteReference(receiver)
	argStruct := api.SetReceiverArgs{user, receiverId, receive}
	args, err := json.Marshal(argStruct)
	helpers.CheckForError(err)

	result := rpc.RPCCall(
		toAddr,
		"SetReceiver",
		args,
	)

	var setReceiverResult api.SetReceiverResult
	err = json.Unmarshal(result, &setReceiverResult)
	helpers.CheckForError(err)

	return setReceiverResult.Result
}

var toAddr transport.NetAddr
