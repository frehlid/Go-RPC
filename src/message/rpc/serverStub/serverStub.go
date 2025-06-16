package serverStub

import (
	"encoding/json"

	"local/lib/helpers"
	"local/message"
	"local/message/rpc/api"
	"local/message/rpc/proxy"
)

func Send(args []byte) []byte {
	var sendArgs api.SendArgs
	err := json.Unmarshal(args, &sendArgs)
	helpers.CheckForError(err)

	result := message.Send(
		sendArgs.User,
		sendArgs.To,
		sendArgs.Text,
	)

	resultStruct := api.SendResult{result}
	resultBuf, err := json.Marshal(resultStruct)
	helpers.CheckForError(err)
	return resultBuf
}

func Receive(args []byte) []byte {
	var receiveArgs api.ReceiveArgs
	err := json.Unmarshal(args, &receiveArgs)
	helpers.CheckForError(err)

	result := message.Receive(
		receiveArgs.User,
	)

	resultStruct := api.ReceiveResult{result}
	resultBuf, err := json.Marshal(resultStruct)
	helpers.CheckForError(err)
	return resultBuf
}

func SetSendingAllowed(args []byte) []byte {
	var setSendingAllowedArgs api.SetSendingAllowedArgs
	err := json.Unmarshal(args, &setSendingAllowedArgs)
	helpers.CheckForError(err)

	message.SetSendingAllowed(
		setSendingAllowedArgs.User,
		setSendingAllowedArgs.From,
		setSendingAllowedArgs.Allowed,
	)

	resultStruct := api.SetSendingAllowedResult{Result: true}
	resultBuf, err := json.Marshal(resultStruct)
	helpers.CheckForError(err)
	return resultBuf
}

func SetReceiver(args []byte) []byte {
	var setReceiverArgs api.SetReceiverArgs

	err := json.Unmarshal(args, &setReceiverArgs)
	helpers.CheckForError(err)

	receiverProxy := proxy.RemoteObjectIdToProxy(setReceiverArgs.Receiver)

	helpers.VerbosePrint("SETTING RECEIVER: ", receiverProxy)
	helpers.VerbosePrint("RECEIVE: ", setReceiverArgs.Receive)
	message.SetReceiver(
		setReceiverArgs.User,
		receiverProxy,
		setReceiverArgs.Receive,
	)

	resultStruct := api.SetReceiverResult{Result: true}
	resultBuf, err := json.Marshal(resultStruct)
	helpers.CheckForError(err)
	return resultBuf
}
