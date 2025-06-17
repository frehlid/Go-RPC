## Run Instructions
There's a small message client built on top of the RPC service. To run it:

1. Start the db daemon:
   `go run cmd/dbd/dbd.go &`
   - It will start listening on an address/port and print out that address.
2. Start the message daemon:
   `go run cmd/messaged/messaged.go addr:port &`
   - where `addr:port` is the address/port the db daemon is listening on.
3. Start the auth daemon:
   `go run cmd/authd/authd.go addr:port &`
4. Finally, connect as many chat clients as you would like.
   `go run cmd/chat/chat.go 127.0.0.1:46756 s/l user pass`
   -  `s` will sign a new user up, `l` will allow you to login to an existing user account.
   - To send a message, use:
        - `s <user> <message>`
   - To read your inbox, use:
        -  `read`
   -  To allow another user to send you messages, use:
        - `allow <user>`
   - To block another user, use:
        - `block <user>`

Clearly and concisely addresses the following points.
1. Evaluate your *rpc* solution from the perspective of transparency.
   - **Access**: 
     - Clients do not need to know if a procedure is executed locally or remotely. The marshalling and transport logic are handled by the client stubs and therefore when the client calls any service, it can use the stub as if it were executing the function locally. 
     - Similarly, the servers (daemons) don't care if they are receiving a request to execute a function locally or remotely, abstracting away marshalling and transporting logic to the server stubs.
     - These points are evident from the fact that the pre-existing services have not been modified in any way other than changing the imports.
   - **Migration**:
     - Our RPC system also permits transparency from the perspective of migration. If we wanted to move a service to a new host, this would be simple as the services register their address into the nameserver database (dbd).
     - Migration would simply look like starting up a new instance of the service, and having the clients restart to re-query the nameserver for the new service address. 
   - **Concurrency**:
     - Our RPC system also allows multiple users to access the same services from any number of devices. In particular, a client can log into the chat system from multiple devices concurrently, and receive messages on both. For the proxy solution, each client is permitted to register one reciever per host. This means that if a client is logged into the service from two different devices, they can choose if they want to have messages pushed on a per-device basis. For example, one could log in from device A and device B, and then choose to have messages pushed to only one device or both. One limitation is that each host is only permitted to register one receiver.
2. Evaluate it from the perspective of procedural modularity (e.g., how much code is duplicated in the client and server stubs)
   - We have made some efforts to reduce duplication within our client and server stubs. For example, the logic of calling an RPC method is shared between all client stubs and is defined in `lib/rpc/rpc.go`. Likewise, the handler method which is responsible for dispatching a received call to the corresponding server stub is shared amongst all server stubs.
   - There is some duplication in the client and server stubs since the process of marshalling arguments is quite repetitive; however, each stub needs to handle marshalling/unmarshalling it's specific arguments at some point. We have tried to abstract this logic as much as possible by defining methods to marshall and unmarshal arguments into a shared RPCData struct, and then having each stub unmarshal the data into method-specific structs defined in the `api` files. This is a bit complex, and requires marshalling/unmarshalling twice (once into RPCData and then into argument specific struct), but improves the modularity of the system since there is only one handler method.
3. Evaluate your `transport` package for trade-offs/limitations
   - Our `transport` package is comprised of a synchronous `Call` function and an asynchronous `Listen` function with one thread listening on a specific port, accepting incoming connections to that port, and dispatching one thread per client connection to handle reading and replying back to that client. This logic was extracted into a seperate function (`handleServerConnection`).
   - One limitation of our `transport` package is the overhead of TCP. The process of setting up a TCP connection is quite cumbersome, and our transport logic does this frequently. We have tried to mitigate this overhead by keeping the TCP connections alive, and reusing them in subsequent calls to the same service if possible. If a client encounters a connection that has broken or been closed by the server, it will then try to re-open a connection to the server and keep that connection for later use. The benifit to this tradeoff is that TCP provides reliability, ensuring messages are received.
   - Another limitaiton has to do with the use of a thread for each connection we are listening to. This potentially poses a scalability issue, since if we have a large number of clients who are communicating with the service infrequently, we are potentially wasting resources keeping these threads alive. Furthermore, there is some overhead when the cpu is switching context between different threads, which is wasteful because this is time when the CPU is not doing useful work.
   - Our `Call` function is synchronous, and while this is not an issue for the `chat` application, it is a limitation of our `transport` package. If a client wanted to make multiple calls concurrently with our existing solution they may encounter errors, especially since the TCP connections are reused.
   - Additionally, we changed the function signature of `Listen` to take in an additional `funcMap` parameter and a generalized `Handler` function which allowed us to abstract all logic related to handling communication between the client and server stubs into one function. This logic was previously duplicated across all of our daemons, so we made this modification to improve the modularity of our system. The trade-off of this decision is increased complexity of communication between the `transport` package and the daemons as each daemon needs to pass in a structure to `Listen` which maps string function names to the serverStub implementations of those functions. Additionally, having an abstracted handler restricts any individual daemon from making custom modifications to its handling logic.
   - Our `transport` also defines a constant `MTU` size of `2KB` for each buffer sent through the network. In practice, this may be too small for very large messages.
