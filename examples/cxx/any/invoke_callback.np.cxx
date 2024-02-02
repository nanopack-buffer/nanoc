// AUTOMATICALLY GENERATED BY NANOPACK. DO NOT MODIFY BY HAND.

#include <nanopack/reader.hxx>
#include <nanopack/writer.hxx>

#include "invoke_callback.np.hxx"

InvokeCallback::InvokeCallback(int32_t handle, NanoPack::Any args)
    : handle(handle), args(std::move(args)) {}

InvokeCallback::InvokeCallback(const NanoPack::Reader &reader,
                               int &bytes_read) {
  const auto begin = reader.begin();
  int ptr = 12;

  const int32_t handle = reader.read_int32(ptr);
  ptr += 4;
  this->handle = handle;

  const int32_t args_byte_size = reader.read_field_size(1);
  args = NanoPack::Any(begin + ptr, begin + ptr + args_byte_size);
  ptr += args_byte_size;

  bytes_read = ptr;
}

InvokeCallback::InvokeCallback(std::vector<uint8_t>::const_iterator begin,
                               int &bytes_read)
    : InvokeCallback(NanoPack::Reader(begin), bytes_read) {}

std::vector<uint8_t> InvokeCallback::data() const {
  std::vector<uint8_t> buf(12);
  NanoPack::Writer writer(&buf);

  writer.write_type_id(TYPE_ID);

  writer.write_field_size(0, 4);
  writer.append_int32(handle);

  writer.write_field_size(1, args.size());
  writer.append_bytes(args.data());

  return buf;
}
