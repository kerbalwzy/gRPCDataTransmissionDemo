package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"

	demo "../protos"
)

const (
	ServerAddress = "127.0.0.1:23333"
	ServerId      = 0
)

func main() {
	listener, err := net.Listen("tcp", ServerAddress)
	if nil != err {
		log.Fatal(err)
	}

	server := grpc.NewServer()
	demo.RegisterGRPCDemoServer(server, &DemoServer{Id: ServerId})

	log.Println("------------------start Golang gRPC server")
	err = server.Serve(listener)
	if nil != err {
		log.Fatal(err)
	}
}

// 随便创建一个结构体，并实现 demo.GRPCDemoServer 接口
// create a struct and implement the `demo.GRPCDemoServer` interface
type DemoServer struct {
	Id int64
	wt sync.WaitGroup
}

// 一元模式(在一次调用中, 客户端只能向服务器传输一次请求数据, 服务器也只能返回一次响应)
// unary-unary(In a single call, the client can only send request once, and the server can
// only respond once.)
func (obj *DemoServer) SimpleMethod(ctx context.Context, request *demo.Request) (*demo.Response, error) {
	log.Printf("SimpleMethod call by Client: Cid = %d Message = %s", request.ClientId, request.RequestData)
	response := &demo.Response{ServerId: obj.Id, ResponseData: "Go server SimpleMethod Ok!!!!"}
	return response, nil
}

// 客户端流模式（在一次调用中, 客户端可以多次向服务器传输数据, 但是服务器只能返回一次响应）
// stream-unary (In a single call, the client can transfer data to the server several times,
// but the server can only return a response once.)
func (obj *DemoServer) ClientStreamingMethod(stream demo.GRPCDemo_ClientStreamingMethodServer) error {
	log.Println("ClientStreamingMethod called, begin to get requests from client ... ")
	// 开始不断接收客户端发送来端数据
	for {
		request, err := stream.Recv()
		if io.EOF == err {
			log.Println("recv done")
			response := &demo.Response{ServerId: ServerId, ResponseData: "Go server ClientStreamingMethod OK!!!!"}
			return stream.SendAndClose(response) // 事实上client只能接收到resp不能接收到return返回端error
		}
		if nil != err {
			log.Println(err)
			return err
		}
		log.Printf("recv from client(%d) message : %s", request.ClientId, request.RequestData)
	}
}

// 服务端流模式（在一次调用中, 客户端只能一次向服务器传输数据, 但是服务器可以多次返回响应）
// unary-stream (In a single call, the client can only transmit data to the server at one time,
// but the server can return the response many times.)
func (obj *DemoServer) ServerStreamingMethod(request *demo.Request, stream demo.GRPCDemo_ServerStreamingMethodServer) error {
	clientId, clientMsg := request.ClientId, request.RequestData
	log.Printf("ServerStreamingMethod called by client(%d) init message:= %s", clientId, clientMsg)

	// 开始多次向client发送数据
	// send data to the client several times
	log.Println("ServerStreamingMethod begin to send message to client...")
	for i := 0; i < 5; i++ {
		msg := fmt.Sprintf("Go server ServerStreamingMethod (%d) OK!!!!", i)
		response := &demo.Response{ServerId: obj.Id, ResponseData: msg}
		if err := stream.Send(response); nil != err {
			log.Println(err)
		}
	}
	log.Printf("send message to client(%d) over", clientId)

	return nil
}

// 双向流模式 (在一次调用中, 客户端和服务器都可以向对方多次收发数据)
// stream-stream (In a single call, both client and server can send and receive data
// to each other multiple times.)
func (obj *DemoServer) BidirectionalStreamingMethod(stream demo.GRPCDemo_BidirectionalStreamingMethodServer) error {

	obj.wt.Add(1)
	// 接收数据
	// recv data
	go func() {
		defer obj.wt.Done()
		log.Println("BidirectionalStreamingMethod called, begin to recv data from client ... ")
		for {
			request, err := stream.Recv()
			if io.EOF == err {
				log.Println("recv from client over")
				break
			}
			if nil != err {
				log.Println(err)
			}
			log.Printf("recve from client(%d) message: %s", request.ClientId, request.RequestData)

		}

	}()

	// 开始多次向client发送数据
	// send data to the server several times
	log.Println("BidirectionalStreamingMethod begin to send message to client...")
	for i := 0; i < 5; i++ {
		msg := fmt.Sprintf("Go server BidirectionalStreamingMethod (%d) OK!!!!", i)
		response := &demo.Response{ServerId: obj.Id, ResponseData: msg}
		if err := stream.Send(response); nil != err {
			log.Println(err)
		}
	}

	obj.wt.Wait()
	return nil
}
