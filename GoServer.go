package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"

	demo "./customGrpcPackages"
)

const (
	ServerAddress = "127.0.0.1:23333"
	ServerId      = 0
)

func main() {
	// 开启GRPC的服务并监听请求
	listener, err := net.Listen("tcp", ServerAddress)
	if nil != err {
		log.Fatal(err)
	}

	server := grpc.NewServer()
	demo.RegisterGRPCDemoServer(server, &DemoServer{Id: ServerId})
	log.Println("------------------start Golang GRPC server")
	err = server.Serve(listener)
	if nil != err {
		log.Fatal(err)
	}
}

// 随便创建一个结构体，并实现 demo.GRPCDemoServer 接口
type DemoServer struct {
	Id int64
	wt sync.WaitGroup
}

// 简单模式
func (obj *DemoServer) SimpleMethod(ctx context.Context, req *demo.Request) (*demo.Response, error) {
	log.Printf("SimpleMethod call by Client: Cid = %d Message = %s", req.Cid, req.ReqMsg)
	resp := &demo.Response{Sid: obj.Id, RespMsg: "Go server SimpleMethod Ok!!!!"}
	return resp, nil
}

// 客户端流模式（在一次调用中, 客户端可以多次向服务器传输数据, 但是服务器只能返回一次响应）
func (obj *DemoServer) CStreamMethod(stream demo.GRPCDemo_CStreamMethodServer) error {
	log.Println("CStreamMethod called, begin to get requests from client ... ")
	// 开始不断接收客户端发送来端数据
	for {
		req, err := stream.Recv()
		if io.EOF == err {
			log.Println("recv done")
			resp := &demo.Response{Sid: ServerId, RespMsg: "Go server CStreamMethod OK!!!!"}
			return stream.SendAndClose(resp) // 事实上client只能接收到resp不能接收到return返回端error
		}
		if nil != err {
			log.Println(err)
			return err
		}
		log.Printf("recv from client(%d) message : %s", req.Cid, req.ReqMsg)
	}
}

// 服务端流模式（在一次调用中, 客户端只能一次向服务器传输数据, 但是服务器可以多次返回响应）
func (obj *DemoServer) SStreamMethod(req *demo.Request, stream demo.GRPCDemo_SStreamMethodServer) error {
	clientId, clientMsg := req.Cid, req.ReqMsg
	log.Printf("SStreamMethod called by client(%d) init message:= %s", clientId, clientMsg)

	// 开始多次向client发送数据
	log.Println("begin to send message to client...")
	for i := 0; i < 5; i++ {
		msg := fmt.Sprintf("Go server SStreamMethod (%d) OK!!!!", i)
		resp := &demo.Response{Sid: obj.Id, RespMsg: msg}
		if err := stream.Send(resp); nil != err {
			log.Println(err)
		}
	}
	log.Printf("send message to client(%d) over", clientId)

	return nil
}

// 双向流模式 (在一次调用中, 客户端和服务器都可以向对象多次收发数据)
func (obj *DemoServer) TWFMethod(stream demo.GRPCDemo_TWFMethodServer) error {
	// 并发收发数据
	obj.wt.Add(1)
	// 接收数据
	go func() {
		defer obj.wt.Done()
		log.Println("TWFMethod called, begin to recv data from client ... ")
		for {
			req, err := stream.Recv()
			if io.EOF == err {
				log.Println("recv from client over")
				break
			}
			if nil != err {
				log.Println(err)
			}
			log.Printf("recve from client(%d) message: %s", req.Cid, req.ReqMsg)

		}

	}()

	// 开始多次向client发送数据
	log.Println("TWFMethod begin to send message to client...")
	for i := 0; i < 5; i++ {
		msg := fmt.Sprintf("Go server TWFMethod (%d) OK!!!!", i)
		resp := &demo.Response{Sid: obj.Id, RespMsg: msg}
		if err := stream.Send(resp); nil != err {
			log.Println(err)
		}
	}

	obj.wt.Wait()
	return nil
}
