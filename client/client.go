package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"

	demo "../protos"
)

const (
	ServerAddress = "127.0.0.1:23333"
	ClientId   = 0
)

func main() {

	conn, err := grpc.Dial(ServerAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err.Error())
	}
	defer conn.Close()

	client := demo.NewGRPCDemoClient(conn)

	CallSimpleMethod(client)

	CallClientStreamingMethod(client)

	CallServerStreamingMethod(client)

	CallBidirectionalStreamingMethod(client)
}

// 一元模式(在一次调用中, 客户端只能向服务器传输一次请求数据, 服务器也只能返回一次响应)
// unary-unary(In a single call, the client can only send request once, and the server can
// only respond once.)
func CallSimpleMethod(client demo.GRPCDemoClient) {
	log.Println("--------------Call SimpleMethod Begin---------------")
	request := demo.Request{ClientId: ClientId, RequestData: "SimpleMethod called by Golang client"}

	// context是用来保存上下文的, 比如我们可以在里面设置这个调用的请求超时时间
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	defer cancel()
	response, err := client.SimpleMethod(ctx, &request)
	if nil != err {
		log.Print("Call SimpleMethod error:", err.Error())
		return
	}
	log.Printf("Get SimpleMethod Response: Sid = %d, RespMsg= %s", response.ServerId, response.ResponseData)
	log.Println("--------------Call SimpleMethod Over----------------")
}

// 客户端流模式（在一次调用中, 客户端可以多次向服务器传输数据, 但是服务器只能返回一次响应）
// stream-unary (In a single call, the client can transfer data to the server several times,
// but the server can only return a response once.)
func CallClientStreamingMethod(client demo.GRPCDemoClient) {
	log.Println("--------------Call ClientStreamingMethod Begin--------------")
	// 获取流信息传输对象stream
	// get the data transmission worker `stream`
	stream, err := client.ClientStreamingMethod(context.Background())
	if nil != err {
		log.Fatal(err)
	}

	// 连续向server发送5次信息
	// send data to server 5 times
	request := &demo.Request{ClientId: ClientId}
	for i := 0; i < 5; i++ {
		request.RequestData = fmt.Sprintf("Golang Client stream message(%d)", i)
		err := stream.Send(request)
		if nil != err {
			log.Println(err)
		}
	}

	// 关闭发送, 并等待接收服务器的响应, 接收完成退出循环
	// close the sender, and wait the server's response data
	for {
		response, err := stream.CloseAndRecv()
		if io.EOF == err {
			log.Println("recv data from server done")
			goto OVER
		}
		if nil != err {
			log.Println(err)
		}
		log.Printf("recv data from server(%d) massege: %s", response.ServerId, response.ResponseData)
	}

OVER:
	log.Println("--------------Call ClientStreamingMethod Over---------------")
}

// 服务端流模式（在一次调用中, 客户端只能一次向服务器传输数据, 但是服务器可以多次返回响应）
// unary-stream (In a single call, the client can only transmit data to the server at one time,
// but the server can return the response many times.)
func CallServerStreamingMethod(client demo.GRPCDemoClient) {
	log.Println("--------------Call ServerStreamingMethod Begin--------------")
	// 调用的同时就给server发送了第一次(仅一次)数据
	// send data to server and get the data transmission worker `stream`
	request := &demo.Request{ClientId: ClientId, RequestData: "SStreamMethod called by Golang client"}
	stream, err := client.ServerStreamingMethod(context.Background(), request)
	if nil != err {
		log.Fatal(err)
	}
	// 不断尝试接收server返回的数据
	for {
		response, err := stream.Recv()
		if io.EOF == err {
			log.Println("recv done")
			goto OVER
		}
		if nil != err {
			log.Println(err)
		}
		log.Printf("recv from server(%d) message : %s", response.ServerId, response.ResponseData)
	}

OVER:
	log.Println("--------------Call ServerStreamingMethod Over---------------")
}

// 双向流模式 (在一次调用中, 客户端和服务器都可以向对方多次收发数据)
// stream-stream (In a single call, both client and server can send and receive data
// to each other multiple times.)
func CallBidirectionalStreamingMethod(client demo.GRPCDemoClient) {
	log.Println("--------------Call BidirectionalStreamingMethod Begin---------------")
	stream, err := client.BidirectionalStreamingMethod(context.Background())
	if nil != err {
		log.Fatal(err)
	}

	wt := new(sync.WaitGroup)
	go func() {
		// 接收数据
		// recv data
		wt.Add(1)
		defer wt.Done()
		for {
			response, err := stream.Recv()
			if io.EOF == err {
				log.Println("recv from server done")
				break
			}
			if nil != err {
				log.Println(err)
			}
			log.Printf("recv from server(%d) message: %s", response.ServerId, response.ResponseData)
		}
	}()

	// 发送数据
	// send data
	log.Println("begin to send message to server...")
	for i := 0; i < 5; i++ {
		request := &demo.Request{ClientId: ClientId, RequestData: fmt.Sprintf("Golang client message(%d)", i)}
		err := stream.Send(request)
		if nil != err {
			log.Println(err)
		}
	}

	_ = stream.CloseSend()
	wt.Wait()

	log.Println("--------------Call BidirectionalStreamingMethod Over---------------")
}
