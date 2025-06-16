package api

import (
	authTypes "local/auth/types"
	"local/message/rpc/proxy"
	messageTypes "local/message/types"
)

type SendArgs struct {
	User authTypes.UserCap
	To   string
	Text string
}

type SendResult struct {
	Result bool
}

type ReceiveArgs struct {
	User authTypes.UserCap
}

type ReceiveResult struct {
	Message *messageTypes.Message
}

type SetSendingAllowedArgs struct {
	User    authTypes.UserCap
	From    string
	Allowed bool
}

type SetSendingAllowedResult struct {
	Result bool
}

type SetReceiverArgs struct {
	User     authTypes.UserCap
	Receiver *proxy.RemoteObjectId
	Receive  bool
}

type SetReceiverResult struct {
	Result bool
}
