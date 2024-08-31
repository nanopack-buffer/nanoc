#include <exception>
#include <nanopack/reader.hxx>
#include <nanopack/writer.hxx>
#include <string>

#include "example_service.np.hxx"

ExampleServiceServer::ExampleServiceServer(NanoPack::RpcServerChannel &channel)
    : NanoPack::RpcServer(channel), handlers(2) {
  handlers.emplace("add", &ExampleServiceServer::add);
  handlers.emplace("subtract", &ExampleServiceServer::subtract);
}

NanoPack::RpcServer::MethodCallResult
ExampleServiceServer::on_method_call(const std::string_view &method,
                                     uint8_t *request_data, size_t offset,
                                     NanoPack::MessageId msg_id) {
  const auto handler = handlers.find(method);
  if (handler == handlers.end()) {
    throw std::invalid_argument("Unknown method " + std::string(method) +
                                " called on ExampleService.");
  }
  return (this->*(handler->second))(request_data, offset, msg_id);
}
NanoPack::RpcServer::MethodCallResult
ExampleServiceServer::_add(uint8_t *request_data, size_t offset,
                           NanoPack::MessageId msg_id) {
  NanoPack::Reader reader(request_data);
  size_t ptr = offset;
  int32_t a;
  reader.read_int32(ptr, a);
  ptr += 4;
  int32_t b;
  reader.read_int32(ptr, b);
  ptr += 4;

  int32_t result = add(a, b);
  NanoPack::Writer writer(6 + 4);
  writer.append_uint8(NanoPack::RpcMessageType::Response);
  writer.append_uint32(msg_id);
  writer.append_uint8(0);
  writer.append_int32(result);
  return {writer.into_data(), writer.size()};
}

NanoPack::RpcServer::MethodCallResult
ExampleServiceServer::_subtract(uint8_t *request_data, size_t offset,
                                NanoPack::MessageId msg_id) {
  NanoPack::Reader reader(request_data);
  size_t ptr = offset;
  int32_t a;
  reader.read_int32(ptr, a);
  ptr += 4;
  int32_t b;
  reader.read_int32(ptr, b);
  ptr += 4;

  int32_t result = subtract(a, b);
  NanoPack::Writer writer(6 + 4);
  writer.append_uint8(NanoPack::RpcMessageType::Response);
  writer.append_uint32(msg_id);
  writer.append_uint8(0);
  writer.append_int32(result);
  return {writer.into_data(), writer.size()};
}

std::future<int32_t> ExampleServiceClient::add(int32_t a, int32_t b) {
  NanoPack::Writer writer(9 + 3 + 8);
  const auto msg_id = new_message_id();
  writer.append_uint8(NanoPack::RpcMessageType::Request);
  writer.append_uint32(msg_id);
  writer.append_uint32(3);
  writer.append_string_view("add");
  writer.append_int32(a);
  writer.append_int32(b);

  return std::async(
      [this](uint32_t msg_id, uint8_t *req_data, size_t req_size) {
        auto res_data =
            send_request_data_async(msg_id, req_data, req_size).get();
        NanoPack::Reader reader(res_data);
        size_t ptr = 0;
        uint8_t err_flag;
        reader.read_uint8(ptr++, err_flag);
        if (err_flag == 1) {
          throw std::runtime_error("RPC on Example::add failed.");
        }
        free(req_data);

        int32_t result;
        reader.read_int32(ptr, result);
        ptr += 4;
        return result;
      },
      msg_id, writer.into_data(), writer.size());
}

std::future<int32_t> ExampleServiceClient::subtract(int32_t a, int32_t b) {
  NanoPack::Writer writer(9 + 8 + 8);
  const auto msg_id = new_message_id();
  writer.append_uint8(NanoPack::RpcMessageType::Request);
  writer.append_uint32(msg_id);
  writer.append_uint32(8);
  writer.append_string_view("subtract");
  writer.append_int32(a);
  writer.append_int32(b);

  return std::async(
      [this](uint32_t msg_id, uint8_t *req_data, size_t req_size) {
        auto res_data =
            send_request_data_async(msg_id, req_data, req_size).get();
        NanoPack::Reader reader(res_data);
        size_t ptr = 0;
        uint8_t err_flag;
        reader.read_uint8(ptr++, err_flag);
        if (err_flag == 1) {
          throw std::runtime_error("RPC on Example::subtract failed.");
        }
        free(req_data);

        int32_t result;
        reader.read_int32(ptr, result);
        ptr += 4;
        return result;
      },
      msg_id, writer.into_data(), writer.size());
}
