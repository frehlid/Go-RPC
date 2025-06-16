package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	auth "local/auth/rpc/clientStub"
	authTypes "local/auth/types"
	"local/cmd/chat/listener"
	db "local/db/rpc/clientStub"
	"local/lib/finalizer"
	"local/lib/handler"
	"local/lib/transport"
	message "local/message/rpc/clientStub"
	"local/message/rpc/proxy"
	"local/message/types"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 5 {
		fmt.Println("usage: address chat s|l user password")
		return
	}
	dbServerAddr := os.Args[1]
	signupOrLogin := os.Args[2]
	userId := os.Args[3]
	password := os.Args[4]

	var cap authTypes.UserCap

	// Bind to db server and retrieve auth/message server addresses
	db.Bind(transport.NewNetAddr(dbServerAddr))
	authAddr := db.Get("auth")
	messageAddr := db.Get("message")

	// bind message and auth to the servers
	auth.Bind(transport.NewNetAddr(authAddr))
	message.Bind(transport.NewNetAddr(messageAddr))

	// Setup a finalizer context to avoid "memory leak" of receiver ref on exit
	ctx, cancel := finalizer.WithCancel(context.Background())
	defer func() { cancel(); <-ctx.Done() }()
	finalizer.AfterFunc(ctx, func() {
		if cap != 0 {
			message.SetReceiver(cap, listener.MessageListener(cap), false)
		}
	})

	go transport.Listen(ctx, proxy.ChatFuncMap, handler.Handler)

	if signupOrLogin == "s" {
		if !auth.Signup(userId, password) {
			fmt.Println("signup failure")
			return
		}
	}
	cap = auth.Login(userId, password)
	if cap == 0 {
		fmt.Println("login failure")
		return
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			fmt.Print("> ")
			cmd, err := reader.ReadString('\n')
			select {
			case <-ctx.Done():
				return
			default:
				if err == io.EOF {
					return
				}
				if err != nil {
					panic(err)
				}
				toks := strings.Fields(string(cmd))
				if len(toks) == 0 {
					fmt.Println(usage)
					continue
				}
				switch toks[0] {
				case "allow":
					message.SetSendingAllowed(cap, toks[1], true)
					fmt.Printf("User %s allowed to send messages\n", toks[1])
				case "block":
					message.SetSendingAllowed(cap, toks[1], false)
					fmt.Printf("User %s blocked from sending messages\n", toks[1])
				case "s":
					sent := message.Send(cap, toks[1], strings.Join(toks[2:], " "))
					if !sent {
						fmt.Println("Invalid receiver")
					}
				case "read":
					readAll(cap)
				case "push":
					fmt.Println("New messages will be pushed automatically")
					message.SetReceiver(cap, listener.MessageListener(cap), true)
					readAll(cap)
				case "pull":
					message.SetReceiver(cap, listener.MessageListener(cap), false)
					fmt.Println("You must use 'read' to check for new messages")
				default:
					fmt.Println(usage)
				}
			}
		}
	}
}

func readAll(cap authTypes.UserCap) {
	for {
		msg := message.Receive(cap)
		if msg == nil {
			break
		}
		printMsg(msg)
	}
}

var usage = "usage: s <user> <message> | read | push | pull | allow <user> | block <user>"

func printMsg(msg *types.Message) {
	fmt.Printf("%s: %s\n", msg.From, msg.Text)
}
