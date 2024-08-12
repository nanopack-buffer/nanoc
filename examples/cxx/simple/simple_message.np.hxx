// AUTOMATICALLY GENERATED BY NANOC

#ifndef SIMPLE_MESSAGE_NP_HXX
#define SIMPLE_MESSAGE_NP_HXX
#include <memory>
#include <nanopack/message.hxx>
#include <nanopack/nanopack.hxx>
#include <nanopack/reader.hxx>
#include <optional>
#include <string>
#include <unordered_map>
#include <vector>

struct SimpleMessage : NanoPack::Message {
  static constexpr NanoPack::TypeId TYPE_ID = 3338766369;

  std::string string_field;
  int32_t int_field;
  double double_field;
  std::optional<std::string> optional_field;
  std::vector<uint8_t> array_field;
  std::unordered_map<std::string, bool> map_field;
  std::unique_ptr<NanoPack::Message> any_message;

  SimpleMessage() = default;

  SimpleMessage(std::string string_field, int32_t int_field,
                double double_field, std::optional<std::string> optional_field,
                std::vector<uint8_t> array_field,
                std::unordered_map<std::string, bool> map_field,
                std::unique_ptr<NanoPack::Message> any_message);

  size_t read_from(NanoPack::Reader &reader);

  [[nodiscard]] NanoPack::Message &get_any_message() const;

  size_t write_to(NanoPack::Writer &writer, int offset) const override;

  [[nodiscard]] NanoPack::TypeId type_id() const override;

  [[nodiscard]] size_t header_size() const override;
};

#endif