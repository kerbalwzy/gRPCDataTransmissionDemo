package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"

	demo "./customGrpcPackages"
)

const (
	GoGrpcServerAddress = "127.0.0.1:23333"
	PyGrpcServerAddress = "127.0.0.1:23334"
	ClientId            = 0
)

func main() {

	conn, err := grpc.Dial(GoGrpcServerAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err.Error())
	}
	defer conn.Close()

	client := demo.NewGRPCDemoClient(conn)

	// client的类型实际上是个指针， 所有下面的方法尽量不要同时调用，实测中同时调用有可能导致服务器出错
	/*
		出错信息如下：
		rpc error: code = Internal desc = transport: transport: the stream is done or WriteHeader was already called
		rpc error: code = Unavailable desc = transport is closing
		TWFMethod called, begin to recv data from client ...
		rpc error: code = Canceled desc = context canceled
		panic: runtime error: invalid memory address or nil pointer dereference
	*/
	// 每个方法单独测试时，没有问题

	//CallSimpleMethod(client)
	//
	//CallCStreamMethod(client)
	//
	//CallSStreamMethod(client)
	//
	CallTWFMethod(client)
}

// 简单模式下，直接调用client的相应方法就是普通的数据传输
func CallSimpleMethod(client demo.GRPCDemoClient) {
	log.Println("--------------Call SimpleMethod Begin---------------")
	req := demo.Request{Cid: ClientId, ReqMsg: "SimpleMethod called by Golang client"}

	// context是用来保存上下文的, 比如我们可以在里面设置这个调用的请求超时时间
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	defer cancel()
	resp, err := client.SimpleMethod(ctx, &req)
	if nil != err {
		log.Print("Call SimpleMethod error:", err.Error())
	}
	log.Printf("Get SimpleMethod Response: Sid = %d, RespMsg= %s", resp.Sid, resp.RespMsg)
	log.Println("--------------Call SimpleMethod Over----------------")
}

// 客户端流模式（在一次调用中, 客户端可以多次向服务器传输数据, 但是服务器只能返回一次响应）
// stream对象只有 Send 和 CloseAndRecv 两种方法。
func CallCStreamMethod(client demo.GRPCDemoClient) {
	log.Println("--------------Call CStreamMethod Begin--------------")
	// 获取流信息传输对象stream
	stream, err := client.CStreamMethod(context.Background())
	if nil != err {
		log.Fatal(err)
	}

	// 连续向server发送5次信息
	req := &demo.Request{Cid: ClientId}
	for i := 0; i < 5; i++ {
		req.ReqMsg = fmt.Sprintf("Golang Client stream message(%d)", i)
		err := stream.Send(req)
		if nil != err {
			log.Println(err)
		}
	}

	// 关闭发送, 并等待接收服务器的响应, 接收完成推出循环
	for {
		resp, err := stream.CloseAndRecv()
		if io.EOF == err {
			log.Println("recv data from server done")
			goto OVER
		}
		if nil != err {
			log.Println(err)
		}
		log.Printf("recv data from server(%d) massege: %s", resp.Sid, resp.RespMsg)
	}

OVER:
	log.Println("--------------Call CStreamMethod Over---------------")
}

// 服务端流模式（在一次调用中, 客户端只能一次向服务器传输数据, 但是服务器可以多次返回响应）
func CallSStreamMethod(client demo.GRPCDemoClient) {
	log.Println("--------------Call SStreamMethod Begin--------------")
	// 调用的同时就给server发送了第一次(仅一次)数据
	req := &demo.Request{Cid: ClientId, ReqMsg: "SStreamMethod called by Golang client"}
	stream, err := client.SStreamMethod(context.Background(), req)
	if nil != err {
		log.Fatal(err)
	}
	// 不断尝试接收server返回的数据
	for {
		resp, err := stream.Recv()
		if io.EOF == err {
			log.Println("recv done")
			goto OVER
		}
		if nil != err {
			log.Println(err)
		}
		log.Printf("recv from server(%d) message : %s", resp.Sid, resp.RespMsg)
	}

OVER:
	log.Println("--------------Call SStreamMethod Over---------------")
}

// 双向流模式 (在一次调用中, 客户端和服务器都可以向对象多次收发数据)
func CallTWFMethod(client demo.GRPCDemoClient) {
	log.Println("--------------Call TWFMethod Begin---------------")
	stream, err := client.TWFMethod(context.Background())
	if nil != err {
		log.Fatal(err)
	}
	// 并发收发数据
	wt := new(sync.WaitGroup)
	go func() {
		// 接收数据
		wt.Add(1)
		defer wt.Done()
		for {
			resp, err := stream.Recv()
			if io.EOF == err {
				log.Println("recv from server done")
				break
			}
			if nil != err {
				log.Println(err)
			}
			log.Printf("recv from server(%d) message: %s", resp.Sid, resp.RespMsg)
		}
	}()

	// 发送数据
	log.Println("begin to send message to server...")
	for i := 0; i < 5; i++ {
		req := &demo.Request{Cid: ClientId, ReqMsg: fmt.Sprintf("Golang client message(%d)", i)}
		err := stream.Send(req)
		if nil != err {
			log.Println(err)
		}
	}
	stream.CloseSend()
	wt.Wait()

	log.Println("--------------Call TWFMethod Over---------------")
}
