# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
import grpc

import demo_pb2 as demo__pb2


class GRPCDemoStub(object):
  """服务service是用来gRPC的方法的, 格式固定
  类似于Python中定义一个类, 类似于Golang中定义一个接口
  """

  def __init__(self, channel):
    """Constructor.

    Args:
      channel: A grpc.Channel.
    """
    self.SimpleMethod = channel.unary_unary(
        '/demo.GRPCDemo/SimpleMethod',
        request_serializer=demo__pb2.Response.SerializeToString,
        response_deserializer=demo__pb2.Response.FromString,
        )
    self.CStreamMethod = channel.stream_unary(
        '/demo.GRPCDemo/CStreamMethod',
        request_serializer=demo__pb2.Request.SerializeToString,
        response_deserializer=demo__pb2.Response.FromString,
        )
    self.SStreamMethod = channel.unary_stream(
        '/demo.GRPCDemo/SStreamMethod',
        request_serializer=demo__pb2.Request.SerializeToString,
        response_deserializer=demo__pb2.Response.FromString,
        )
    self.TWFMethod = channel.stream_stream(
        '/demo.GRPCDemo/TWFMethod',
        request_serializer=demo__pb2.Request.SerializeToString,
        response_deserializer=demo__pb2.Response.FromString,
        )


class GRPCDemoServicer(object):
  """服务service是用来gRPC的方法的, 格式固定
  类似于Python中定义一个类, 类似于Golang中定义一个接口
  """

  def SimpleMethod(self, request, context):
    """简单模式
    """
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def CStreamMethod(self, request_iterator, context):
    """客户端流模式（在一次调用中, 客户端可以多次向服务器传输数据, 但是服务器只能返回一次响应）
    """
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def SStreamMethod(self, request, context):
    """服务端流模式（在一次调用中, 客户端只能一次向服务器传输数据, 但是服务器可以多次返回响应）
    """
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def TWFMethod(self, request_iterator, context):
    """双向流模式 (在一次调用中, 客户端和服务器都可以向对象多次收发数据)
    """
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')


def add_GRPCDemoServicer_to_server(servicer, server):
  rpc_method_handlers = {
      'SimpleMethod': grpc.unary_unary_rpc_method_handler(
          servicer.SimpleMethod,
          request_deserializer=demo__pb2.Response.FromString,
          response_serializer=demo__pb2.Response.SerializeToString,
      ),
      'CStreamMethod': grpc.stream_unary_rpc_method_handler(
          servicer.CStreamMethod,
          request_deserializer=demo__pb2.Request.FromString,
          response_serializer=demo__pb2.Response.SerializeToString,
      ),
      'SStreamMethod': grpc.unary_stream_rpc_method_handler(
          servicer.SStreamMethod,
          request_deserializer=demo__pb2.Request.FromString,
          response_serializer=demo__pb2.Response.SerializeToString,
      ),
      'TWFMethod': grpc.stream_stream_rpc_method_handler(
          servicer.TWFMethod,
          request_deserializer=demo__pb2.Request.FromString,
          response_serializer=demo__pb2.Response.SerializeToString,
      ),
  }
  generic_handler = grpc.method_handlers_generic_handler(
      'demo.GRPCDemo', rpc_method_handlers)
  server.add_generic_rpc_handlers((generic_handler,))