import grpc
import time
from customGrpcPackages import demo_pb2, demo_pb2_grpc

GoGrpcServerAddress = "127.0.0.1:23333"
PyGrpcServerAddress = "127.0.0.1:23334"
ClientId = 1


# 简单模式下，直接调用stub的相应方法就是普通的数据传输
def simple_method(stub):
    print("--------------Call SimpleMethod Begin--------------")
    req = demo_pb2.Request(Cid=ClientId, ReqMsg="called by Python client")
    resp = stub.SimpleMethod(req)
    print(f"get resp from server({resp.Sid}), the message: {resp.RespMsg}")
    print("--------------Call SimpleMethod Over---------------")


# 客户端流模式（在一次调用中, 客户端可以多次向服务器传输数据, 但是服务器只能返回一次响应）
def c_stream_method(stub):
    print("--------------Call CStreamMethod Begin--------------")

    # 创建一个生成器
    def req_msgs():
        for i in range(5):
            req = demo_pb2.Request(Cid=ClientId, ReqMsg=f"called by Python client, message: {i}")
            yield req

    resp = stub.CStreamMethod(req_msgs())
    print(f"get resp from server({resp.Sid}), the message: {resp.RespMsg}")
    print("--------------Call CStreamMethod Over---------------")


# 服务端流模式（在一次调用中, 客户端只能一次向服务器传输数据, 但是服务器可以多次返回响应）
def s_stream_method(stub):
    print("--------------Call SStreamMethod Begin--------------")
    req = demo_pb2.Request(Cid=ClientId, ReqMsg="called by Python client")
    resps = stub.SStreamMethod(req)
    for resp in resps:
        print(f"recv from server({resp.Sid}, message={resp.RespMsg})")

    print("--------------Call SStreamMethod Over---------------")


# 双向流模式 (在一次调用中, 客户端和服务器都可以向对象多次收发数据)
def twf_method(stub):
    print("--------------Call TWFMethod Begin---------------")

    # 创建一个req_msgs生成器
    def req_msgs():
        for i in range(5):
            req = demo_pb2.Request(Cid=ClientId, ReqMsg=f"called by Python client, message: {i}")
            yield req
            time.sleep(1)

    resps = stub.TWFMethod(req_msgs())
    for resp in resps:
        print(f"recv from server({resp.Sid}, message={resp.RespMsg})")

    print("--------------Call TWFMethod Over---------------")


def main():
    with grpc.insecure_channel(PyGrpcServerAddress) as channel:
        stub = demo_pb2_grpc.GRPCDemoStub(channel)
        simple_method(stub)
        c_stream_method(stub)
        s_stream_method(stub)
        twf_method(stub)


if __name__ == '__main__':
    main()
