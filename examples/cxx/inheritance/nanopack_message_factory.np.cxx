// AUTOMATICALLY GENERATED BY NANOPACK. DO NOT MODIFY BY HAND.

#include "nanopack_message_factory.np.hxx"

#include "text.np.hxx"
#include "widget.np.hxx"

std::unique_ptr<NanoPack::Message>
make_nanopack_message(NanoPack::Reader &reader, size_t &bytes_read) {
  switch (reader.read_type_id()) {
  case 1676374721: {
    auto ptr = std::make_unique<Widget>();
    bytes_read = ptr->read_from(reader);
    return ptr;
  }
  case 3495336243: {
    auto ptr = std::make_unique<Text>();
    bytes_read = ptr->read_from(reader);
    return ptr;
  }
  default:
    return nullptr;
  }
}
