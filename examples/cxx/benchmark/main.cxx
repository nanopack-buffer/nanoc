#include "person.np.hxx"
#include "person.pb.h"
#include <chrono>
#include <cstdint>
#include <cstdlib>

#define ITER_COUNT 1000000

long long current_time_ms() {
  return std::chrono::duration_cast<std::chrono::milliseconds>(
             std::chrono::system_clock::now().time_since_epoch())
      .count();
}

long long run_nanopack() {
  std::vector<Person> friends(2);
  friends.emplace_back("Evelyn", std::nullopt, "Hartley", 30,
                       std::vector<Person>(0));
  friends.emplace_back("Nolan", "Thomas", "Blackwood", 30,
                       std::vector<Person>(0));
  Person person("John", std::nullopt, "Doe", 40, std::move(friends));

  uint8_t *buf = (uint8_t *)(std::malloc(500 * sizeof(uint8_t)));
  NanoPack::Reader reader(nullptr);
  NanoPack::Writer writer(buf, 500);
  Person p;

  const auto before = current_time_ms();
  for (int i = 0; i < ITER_COUNT; ++i) {
    writer.reset();
    person.write_to(writer, 0);
    reader.buffer = writer.data();
    p.read_from(reader);
  }
  const auto after = current_time_ms();

  free(buf);

  return after - before;
}

long long run_proto() {
  ProtoPerson evelyn;
  evelyn.set_first_name("Evelyn");
  evelyn.set_last_name("Hartley");
  evelyn.set_age(30);

  ProtoPerson nolan;
  nolan.set_first_name("Nolan");
  nolan.set_middle_name("Thomas");
  nolan.set_last_name("Blackwood");
  nolan.set_age(30);

  ProtoPerson person;
  person.set_first_name("John");
  person.set_last_name("Doe");
  person.set_age(40);
  person.mutable_friends()->Add(std::move(evelyn));
  person.mutable_friends()->Add(std::move(nolan));

  ProtoPerson p2;

  const auto before = current_time_ms();
  for (int i = 0; i < ITER_COUNT; ++i) {
    std::string bytes;
    person.SerializeToString(&bytes);
    p2.ParseFromString(bytes);
  }
  const auto after = current_time_ms();
  return after - before;
}

int main() {
  const long long nanopack_time_ms = run_nanopack();
  const long long proto_time_ms = run_proto();

  std::cout << "nanopack took: " << nanopack_time_ms << "ms" << std::endl;
  std::cout << "protobuf took: " << proto_time_ms << "ms" << std::endl;
}
