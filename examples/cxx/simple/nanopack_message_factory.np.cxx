// AUTOMATICALLY GENERATED BY NANOPACK. DO NOT MODIFY BY HAND.

#include "nanopack_message_factory.np.hxx"

#include "nested_message.np.hxx"
#include "simple_message.np.hxx"

std::unique_ptr<NanoPack::Message>
make_nanopack_message(NanoPack::Reader &reader) {
  size_t _;
  return make_nanopack_message(reader, _);
}

std::unique_ptr<NanoPack::Message>
make_nanopack_message(NanoPack::Reader &reader, size_t &bytes_read) {
  switch (reader.read_type_id()) {
  case 2309634176: {
    auto ptr = std::make_unique<NestedMessage>();
    bytes_read = ptr->read_from(reader);
    return ptr;
  }
  case 3338766369: {
    auto ptr = std::make_unique<SimpleMessage>();
    bytes_read = ptr->read_from(reader);
    return ptr;
  }
  default:
    return nullptr;
  }
}