##  Data transmission demo for using gRPC in Python

Four ways of data transmission when gRPC is used in Python.  [Official Guide](<https://grpc.io/docs/guides/concepts/#unary-rpc>)

- #### unary-unary

  In a single call, the client can only send request once, and the server can only respond once.

  `client.go : CallSimpleMethod`

  `server.go : SimpleMethod`

- #### stream-unary

  In a single call, the client can transfer data to the server an arbitrary number of times, but the server can only return a response once.

  `client.go : CallClientStreamingMethod `

  `server.go : ClientStreamingMethodd`

- #### unary-stream

  In a single call, the client can only transmit data to the server at one time, but the server can return the response many times.

  `client.go : CallServerStreamingMethod`

  `server.go : ServerStreamingMethod`

- #### stream-stream

  In a single call, both client and server can send and receive data 
  to each other multiple times.

  `client.go : CallBidirectionalStreamingMethod`

  `server.go : BidirectionalStreamingMethod`

