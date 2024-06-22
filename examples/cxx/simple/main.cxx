#include "person.np.hxx"
#include <cstdint>
#include <iostream>
#include <memory>
#include <optional>

int main() {
  std::cout << "A simple program demonstrating conversion between NanoPack "
               "data and C++ struct."
            << std::endl;

  Person person;
  person.first_name = "John";
  person.middle_name = std::nullopt;
  person.last_name = "Doe";
  person.age = 40;
  person.other_friend = std::make_unique<Person>();
  person.other_friend->first_name = "Tim";
  person.other_friend->last_name = "Cook";
  person.other_friend->age = 50;

  std::cout << "test" << std::endl;

  NanoPack::Writer writer;
  const size_t bytes_written = person.write_to(writer, 0);

  uint8_t *buf = writer.data();

  std::cout << "raw bytes:";
  for (size_t i = 0; i < bytes_written; ++i) {
    std::cout << +buf[i] << " ";
  }
  std::cout << std::endl;
  std::cout << "size of serialized data in bytes: " << bytes_written
            << std::endl
            << "=====================" << std::endl;

  Person person1;
  NanoPack::Reader reader(buf);
  const size_t bytes_read = person1.read_from(reader);
  std::cout << "First name: " << person1.first_name << std::endl;
  std::cout << "Last name: " << person1.last_name << std::endl;
  if (!person1.middle_name.has_value()) {
    std::cout << "This person does not have a middle name." << std::endl;
  }
  std::cout << "Age: " << +person1.age << std::endl;

  std::cout << "His friend:" << std::endl;
  std::cout << "    First name: " << person1.other_friend->first_name
            << std::endl;
  std::cout << "    Last name: " << person1.other_friend->last_name
            << std::endl;
  if (!person1.other_friend->middle_name.has_value()) {
    std::cout << "    This person does not have a middle name." << std::endl;
  }
  std::cout << "    Age: " << person1.other_friend->age << std::endl;

  return 0;
}
