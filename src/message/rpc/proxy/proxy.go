package proxy

import (
	"bytes"
	"encoding/json"
	"local/lib/helpers"
	"local/lib/rpc"
	"local/lib/transport"
	"local/message/types"
	"math/rand"
)

type RemoteObjectId struct {
	Host string
	Gid  GlobalObjectId
}

type RemoteMessage struct {
	Gid     GlobalObjectId
	Message types.Message
}

type GlobalObjectId uint64
type ReceiverProxy RemoteObjectId

var globalObjectIdToLocalObjectMap = make(map[GlobalObjectId]types.Receiver)
var hostToReceiverProxyMap = make(map[string]*ReceiverProxy)

var ChatFuncMap = map[string]func(args []byte) []byte{
	"HandleMessageReceived": HandleMessageReceived,
}

func HandleMessageReceived(args []byte) []byte {
	var remoteMessage RemoteMessage
	err := json.Unmarshal(args, &remoteMessage)
	helpers.CheckForError(err)

	receiver := GlobalObjectIdToLocalObject(remoteMessage.Gid)
	receiver.MessageReceived(remoteMessage.Message)

	buf, err := json.Marshal(true)
	helpers.CheckForError(err)

	//transport.Reply(bytes.NewBuffer(buf), transport.NewNetAddr(remoteMessage.From))

	return buf
}

func (r ReceiverProxy) MessageReceived(msg types.Message) {
	remoteMsg := RemoteMessage{r.Gid, msg}
	buf, err := json.Marshal(remoteMsg)
	helpers.CheckForError(err)

	rpcData := rpc.RPCData{Method: "HandleMessageReceived", Data: buf}
	remoteBuf, err := json.Marshal(rpcData)
	helpers.CheckForError(err)

	_, err = transport.Call(bytes.NewBuffer(remoteBuf), transport.NewNetAddr(r.Host))
	helpers.CheckForError(err)
}

func LocalObjectToRemoteReference(r types.Receiver) *RemoteObjectId {
	rid := &RemoteObjectId{
		Host: transport.LocalAddr(),
		Gid:  GlobalObjectId(rand.Uint64()),
	}

	globalObjectIdToLocalObjectMap[rid.Gid] = r
	return rid
}

func RemoteObjectIdToProxy(rid *RemoteObjectId) *ReceiverProxy {
	// limitation -> a host can only have one receiver, but a user can have multiple if on different hosts
	p := hostToReceiverProxyMap[rid.Host]

	if p == nil {
		p = &ReceiverProxy{
			Host: rid.Host,
			Gid:  rid.Gid,
		}
		hostToReceiverProxyMap[rid.Host] = p
	}

	return p
}

func GlobalObjectIdToLocalObject(gid GlobalObjectId) types.Receiver {
	return globalObjectIdToLocalObjectMap[gid]
}
