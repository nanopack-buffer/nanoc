#include "nested_message.np.hxx"
#include "simple_message.np.hxx"
#include <cstdint>
#include <iostream>
#include <nanopack/reader.hxx>
#include <optional>
#include <unordered_map>
#include <vector>

int main() {
  std::cout << "A simple program demonstrating conversion between NanoPack "
               "data and C++ struct."
            << std::endl;

  SimpleMessage message("hello world", 123456, 123.456, std::nullopt,
                        std::vector<uint8_t>{1, 2, 3},
                        std::unordered_map<std::string, bool>{{"hello", true}},
                        std::make_unique<NestedMessage>("nested"));

  NanoPack::Writer writer;
  const size_t bytes_written = message.write_to(writer, 0);

  uint8_t *buf = writer.data();

  std::cout << "raw bytes: ";
  for (size_t i = 0; i < bytes_written; ++i) {
    std::cout << +buf[i] << " ";
  }
  std::cout << std::endl;
  std::cout << "size of serialized data in bytes: " << bytes_written
            << std::endl
            << "=====================" << std::endl;

  std::cout << "decoded: ";

  NanoPack::Reader reader(buf);
  SimpleMessage decoded;
  decoded.read_from(reader);

  std::cout << "string_field: " << decoded.string_field << std::endl;

  for (auto &item : decoded.array_field) {
    std::cout << +item << std::endl;
  }
}
