// AUTOMATICALLY GENERATED BY NANOPACK. DO NOT MODIFY BY HAND.

#include <nanopack/reader.hxx>
#include <nanopack/writer.hxx>

#include "click_event.np.hxx"

ClickEvent::ClickEvent(double x, double y, int64_t timestamp)
    : x(x), y(y), timestamp(timestamp) {}

ClickEvent::ClickEvent(const NanoPack::Reader &reader, int &bytes_read) {
  const auto begin = reader.begin();
  int ptr = 16;

  const double x = reader.read_double(ptr);
  ptr += 8;
  this->x = x;

  const double y = reader.read_double(ptr);
  ptr += 8;
  this->y = y;

  const int64_t timestamp = reader.read_int64(ptr);
  ptr += 8;
  this->timestamp = timestamp;

  bytes_read = ptr;
}

ClickEvent::ClickEvent(std::vector<uint8_t>::const_iterator begin,
                       int &bytes_read)
    : ClickEvent(NanoPack::Reader(begin), bytes_read) {}

NanoPack::TypeId ClickEvent::type_id() const { return TYPE_ID; }

int ClickEvent::header_size() const { return 16; }

size_t ClickEvent::write_to(std::vector<uint8_t> &buf, int offset) const {
  const size_t buf_size_before = buf.size();

  buf.resize(offset + 16);

  NanoPack::write_type_id(TYPE_ID, offset, buf);

  NanoPack::write_field_size(0, 8, offset, buf);
  NanoPack::append_double(x, buf);

  NanoPack::write_field_size(1, 8, offset, buf);
  NanoPack::append_double(y, buf);

  NanoPack::write_field_size(2, 8, offset, buf);
  NanoPack::append_int64(timestamp, buf);

  return buf.size() - buf_size_before;
}

std::vector<uint8_t> ClickEvent::data() const {
  std::vector<uint8_t> buf(16);
  write_to(buf, 0);
  return buf;
}
