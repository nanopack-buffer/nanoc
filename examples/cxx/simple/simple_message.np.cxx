// AUTOMATICALLY GENERATED BY NANOC

#include <nanopack/reader.hxx>
#include <nanopack/writer.hxx>

#include "nanopack_message_factory.np.hxx"

#include "simple_message.np.hxx"

SimpleMessage::SimpleMessage(std::string string_field, int32_t int_field,
                             double double_field,
                             std::optional<std::string> optional_field,
                             std::vector<uint8_t> array_field,
                             std::unordered_map<std::string, bool> map_field,
                             std::unique_ptr<NanoPack::Message> any_message)
    : string_field(std::move(string_field)), int_field(int_field),
      double_field(double_field), optional_field(std::move(optional_field)),
      array_field(std::move(array_field)), map_field(std::move(map_field)),
      any_message(std::move(any_message)) {}

size_t SimpleMessage::read_from(NanoPack::Reader &reader) {
  uint8_t *buf = reader.buffer;
  int ptr = 32;

  const int32_t string_field_size = reader.read_field_size(0);
  reader.read_string(ptr, string_field_size, string_field);
  ptr += string_field_size;

  reader.read_int32(ptr, int_field);
  ptr += 4;

  reader.read_double(ptr, double_field);
  ptr += 8;

  if (reader.read_field_size(3) < 0) {
    optional_field = std::nullopt;
  } else {
    optional_field = std::move(std::string());
    const int32_t optional_field_size = reader.read_field_size(3);
    reader.read_string(ptr, optional_field_size, optional_field);
    ptr += optional_field_size;
  }

  const int32_t array_field_byte_size = reader.read_field_size(4);
  const int32_t array_field_vec_size = array_field_byte_size / 1;
  array_field.resize(array_field_vec_size);
  for (int i = 0; i < array_field_vec_size; ++i) {
    auto &i_item = array_field[i];
    reader.read_uint8(ptr, i_item);
    ptr += 1;
  }

  uint32_t map_field_map_size;
  reader.read_uint32(ptr, map_field_map_size);
  ptr += 4;
  map_field.reserve(map_field_map_size);
  for (int i = 0; i < map_field_map_size; i++) {
    uint32_t i_key_size;
    reader.read_uint32(ptr, i_key_size);
    ptr += 4;
    std::string i_key;
    reader.read_string(ptr, i_key_size, i_key);
    ptr += i_key_size;
    auto &i_value = map_field[i_key];
    reader.read_bool(ptr++, i_value);
  }

  reader.buffer += ptr;
  size_t any_message_bytes_read;
  any_message =
      std::move(make_nanopack_message(reader, any_message_bytes_read));
  reader.buffer = buf;
  ptr += any_message_bytes_read;

  return ptr;
}

NanoPack::Message &SimpleMessage::get_any_message() const {
  return *any_message;
}

NanoPack::TypeId SimpleMessage::type_id() const { return TYPE_ID; }

size_t SimpleMessage::header_size() const { return 32; }

size_t SimpleMessage::write_to(NanoPack::Writer &writer, int offset) const {
  const size_t writer_size_before = writer.size();

  writer.reserve_header(32);

  writer.write_type_id(TYPE_ID, offset);

  writer.write_field_size(0, string_field.size(), offset);
  writer.append_string(string_field);

  writer.write_field_size(1, 4, offset);
  writer.append_int32(int_field);

  writer.write_field_size(2, 8, offset);
  writer.append_double(double_field);

  if (optional_field.has_value()) {
    writer.write_field_size(3, optional_field->size(), offset);
    writer.append_string(*optional_field);
  } else {
    writer.write_field_size(3, -1, offset);
  }

  const int32_t array_field_byte_size = array_field.size() * 1;
  writer.write_field_size(4, array_field_byte_size, offset);
  for (const auto &i : array_field) {
    writer.append_uint8(i);
  }

  const size_t map_field_map_size = map_field.size();
  writer.append_int32(map_field_map_size);
  int32_t map_field_byte_size = 4 + map_field_map_size * -1;
  for (const auto &i : map_field) {
    auto i_key = i.first;
    auto i_value = i.second;
    writer.append_int32(i_key.size());
    writer.append_string(i_key);
    writer.append_bool(i_value);
    map_field_byte_size += i_key.size() + 4;
  }
  writer.write_field_size(5, map_field_byte_size, offset);

  const size_t any_message_byte_size =
      any_message->write_to(writer, writer.size());
  writer.write_field_size(6, any_message_byte_size, offset);

  return writer.size() - writer_size_before;
}
