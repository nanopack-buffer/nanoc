// AUTOMATICALLY GENERATED BY NANOPACK. DO NOT MODIFY BY HAND.

#include <nanopack/reader.hxx>
#include <nanopack/writer.hxx>

#include "column.np.hxx"

Column::Column(const Alignment &alignment) : alignment(alignment) {}

Column::Column(const NanoPack::Reader &reader, int &bytes_read) {
  const auto begin = reader.begin();
  int ptr = 8;

  const int32_t alignment_size = reader.read_field_size(0);
  const std::string alignment_raw_value =
      reader.read_string(ptr, alignment_size);
  ptr += alignment_size;
  alignment = Alignment(alignment_raw_value);

  bytes_read = ptr;
}

Column::Column(std::vector<uint8_t>::const_iterator begin, int &bytes_read)
    : Column(NanoPack::Reader(begin), bytes_read) {}

std::vector<uint8_t> Column::data() const {
  std::vector<uint8_t> buf(8);
  NanoPack::Writer writer(&buf);

  writer.write_type_id(TYPE_ID);

  writer.write_field_size(0, alignment.value().size());
  writer.append_string(alignment.value());

  return buf;
}
