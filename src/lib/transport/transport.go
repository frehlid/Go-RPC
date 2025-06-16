package transport

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"local/lib/helpers"
	"math/rand"
	"net"
	"os"
	"strings"
	"sync"
)

const MTU = 2048

var listenLocalAddr string = ""

type NetAddr struct {
	address net.Addr
}

// These are to prevent really long function signatures, can't define elsewhere
type MessageHandler func(msg *bytes.Buffer, from NetAddr, funcMap FunctionMap)
type FunctionMap map[string]func(args []byte) []byte

func NewNetAddr(address string) NetAddr {
	addr, _ := net.ResolveTCPAddr("tcp", address)
	return NetAddr{address: addr}
}

var addrConnMapClient = make(map[NetAddr]net.Conn)
var addrConnMapServer = make(map[NetAddr]net.Conn)
var addrConnMapClientMutex sync.Mutex
var addrConnMapServerMutex sync.Mutex

// constructors for RPC connections
func GetConnection(to NetAddr) (net.Conn, error) {
	addrConnMapClientMutex.Lock()
	_, ok := addrConnMapClient[to]
	addrConnMapClientMutex.Unlock()

	if !ok {
		helpers.VerbosePrint("About to connect to " + to.address.String())
		conn, err := net.Dial("tcp", to.address.String())
		helpers.CheckForError(err)
		addrConnMapClient[to] = conn
	}

	addrConnMapClientMutex.Lock()
	conn := addrConnMapClient[to]
	addrConnMapClientMutex.Unlock()

	return conn, nil
}

func GetRandomPort() string {
	return fmt.Sprintf("%d", rand.Intn(65535-1024)+1024)

}

func CloseConnection(conn net.Conn, server bool) {
	err := conn.Close()
	helpers.CheckForError(err)

	if server {
		addrConnMapServerMutex.Lock()
		delete(addrConnMapServer, NetAddr{conn.RemoteAddr()})
		addrConnMapServerMutex.Unlock()
	} else {
		addrConnMapClientMutex.Lock()
		delete(addrConnMapClient, NetAddr{conn.RemoteAddr()})
		addrConnMapClientMutex.Unlock()
	}
}

func Call(payload *bytes.Buffer, to NetAddr) (result *bytes.Buffer, err error) {
	var conn net.Conn

	for {
		conn, err = GetConnection(to)
		helpers.CheckForError(err)

		_, err = conn.Write(payload.Bytes())
		if err != nil {
			// if the error is that the connection is closed/broken, remove it from the map and loop again to get a new connection
			if strings.Contains(err.Error(), "use of closed network connection") || strings.Contains(err.Error(), "broken pipe") || strings.Contains(err.Error(), "connection reset by peer") {
				// remove from map and try again
				addrConnMapClientMutex.Lock()
				delete(addrConnMapClient, NetAddr{conn.RemoteAddr()})
				addrConnMapClientMutex.Unlock()
				continue
			} else {
				helpers.CheckForError(err) // if not one of these errors, panic
			}
		} else {
			break
		}
	}

	rbuf := make([]byte, MTU)
	var n int
	for {
		n, err = conn.Read(rbuf)
		if err == io.EOF {
			helpers.VerbosePrint("Client disconnected")
			CloseConnection(conn, false)
			break
		} else if errors.Is(err, os.ErrDeadlineExceeded) {
			helpers.VerbosePrint("Deadline Exceeded in call")
			continue
		} else {
			helpers.CheckForError(err)
			break // if no issue, we received smthing
		}
	}

	helpers.VerbosePrint("RBUF CALL THREAD: ", string(rbuf[:n]))
	return bytes.NewBuffer(rbuf[:n]), nil
}

func Reply(payload *bytes.Buffer, to NetAddr) {
	addrConnMapServerMutex.Lock()
	conn, ok := addrConnMapServer[to]
	addrConnMapServerMutex.Unlock()

	if !ok {
		panic("Could not find connection to " + to.address.String())
	} else {
		_, writeErr := conn.Write(payload.Bytes())
		helpers.CheckForError(writeErr)
	}
}

func handleServerConnection(ctx context.Context, from NetAddr, conn net.Conn, funcMap FunctionMap, handler MessageHandler) {
	for {
		select {
		case <-ctx.Done():
			helpers.VerbosePrint("ServerConnection closed to ", from.address.String())
			CloseConnection(conn, true)
			return // Exit the function properly
		default:
			rBuf := make([]byte, MTU)
			n, err := conn.Read(rBuf)
			if err == io.EOF {
				helpers.VerbosePrint("Closing connection to " + from.address.String())
				CloseConnection(conn, true)
				return
			} else if errors.Is(err, os.ErrDeadlineExceeded) {
				// do nothing
			} else {
				helpers.CheckForError(err)
			}

			if n > 0 {
				helpers.VerbosePrint("RBUF LISTEN THREAD: ", string(rBuf[:n]))
				go handler(bytes.NewBuffer(rBuf[:n]), from, funcMap)
			}
		}
	}
}

// list of open connections
// Start listenting for incoming calls
func Listen(ctx context.Context, funcMap map[string]func(args []byte) []byte, handler MessageHandler) {
	port := GetRandomPort()

	// specify ip_address:port to listen for incoming connections on
	var listener net.Listener
	var err error
	for {
		serverAddress := "127.0.0.1:" + port
		listener, err = net.Listen("tcp", serverAddress)

		// reliability -- if port is in use, try a different random one...
		if err != nil && strings.Contains(err.Error(), "address already in use") {
			port = GetRandomPort()
			continue
		}
		break
	}

	helpers.CheckForError(err)
	listenLocalAddr = listener.Addr().String()

	// continously loop, accept incoming connections, dispatch a new goroutine to handle each connection
	for {
		select {
		case <-ctx.Done():
			helpers.VerbosePrint("Listen() done")
			return
		default:
			// wait for connection; accept blocks so that we aren't spinning
			conn, err := listener.Accept()
			helpers.CheckForError(err)
			from := NetAddr{conn.RemoteAddr()}

			helpers.VerbosePrint("Accepted connection from: ", from.address.String())

			// store from -> conn for reply (called by handler once it is done)
			addrConnMapServerMutex.Lock()
			addrConnMapServer[from] = conn
			addrConnMapServerMutex.Unlock()

			go handleServerConnection(ctx, from, conn, funcMap, handler)
		}

	}
}

// For the proxy object to access the server
func LocalAddr() string {
	for {
		if listenLocalAddr == "" {
			continue
		}
		return listenLocalAddr
	}
}
