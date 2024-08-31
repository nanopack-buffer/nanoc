#include <future>
#include <nanopack/rpc.hxx>
#include <string_view>
#include <unordered_map>

#ifndef EXAMPLE_SERVICE_NP_HXX
#define EXAMPLE_SERVICE_NP_HXX

class ExampleServiceServer : public NanoPack::RpcServer {
  std::unordered_map<std::string_view,
                     MethodCallResult (ExampleServiceServer::*)(
                         uint8_t *, size_t, NanoPack::MessageId)>
      handlers;

  MethodCallResult on_method_call(const std::string_view &method,
                                  uint8_t *request_data, size_t offset,
                                  NanoPack::MessageId msg_id) override;
  MethodCallResult _add(uint8_t *request_data, size_t offset,
                        NanoPack::MessageId msg_id);
  virtual int32_t add(int32_t a, int32_t b) = 0;

  MethodCallResult _subtract(uint8_t *request_data, size_t offset,
                             NanoPack::MessageId msg_id);
  virtual int32_t subtract(int32_t a, int32_t b) = 0;

public:
  ExampleServiceServer(NanoPack::RpcServerChannel &channel);
};

class ExampleServiceClient : public NanoPack::RpcClient {
public:
  using NanoPack::RpcClient::RpcClient;
  std::future<int32_t> add(int32_t a, int32_t b);

  std::future<int32_t> subtract(int32_t a, int32_t b);
};

#endif
