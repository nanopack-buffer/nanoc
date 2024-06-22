#include <cstddef>
#include <cstdint>
#include <iostream>

#include "column.np.hxx"
#include "nanopack/writer.hxx"
#include "np_date.np.hxx"

int main() {
  std::cout << "A simple program demonstrating enums in NanoPack." << std::endl;
  std::cout << "In NanoPack, enums are serialized into their backing types."
            << std::endl;
  std::cout << "The backing type can be specified in the schema. When none is "
               "specified, nanoc will determine the most appropriate type."
            << std::endl
            << std::endl;

  NanoPack::Writer writer;

  std::cout << "The Date message uses the Week and Month enum, both of which "
               "are backed by an int8."
            << std::endl;

  const NpDate date(19, Week::MONDAY, Month::JUNE, 2000);
  size_t bytes_written = date.write_to(writer, 0);
  uint8_t *date_serialized = writer.data();
  std::cout << "Raw bytes of Date:";
  for (size_t i = 0; i < bytes_written; ++i) {
    std::cout << +date_serialized[i] << " ";
  }
  std::cout << std::endl;
  std::cout << "Total bytes: " << bytes_written << std::endl << std::endl;

  // although not recommended, Writer can be reused,
  // but since Writer::data returns a pointer to its internal buffer,
  // when Writer is used again, any previous serialized data will be
  // overwritten, which means the previously obtained pointers will point to new
  // data instead!
  writer.reset();

  const Column column(Alignment::CENTER);
  bytes_written = column.write_to(writer, 0);
  uint8_t *column_serialized = writer.data();
  std::cout << "The Column message uses the Alignment enum which is backed by "
               "a string."
            << std::endl;
  std::cout << "Raw bytes of Column:";
  for (size_t i = 0; i < bytes_written; ++i) {
    std::cout << +column_serialized[i] << " ";
  }
  std::cout << std::endl;
  std::cout << "Total bytes: " << bytes_written << std::endl << std::endl;
}
