#include <chrono>
#include <cstdint>
#include <iostream>

#include "click_event.np.hxx"
#include "invoke_callback.np.hxx"
#include "nanopack/reader.hxx"
#include "nanopack/writer.hxx"

int main() {
  std::cout << "The any keyword in NanoPack allows you to store any NanoPack "
               "data type."
            << std::endl;

  const auto current_time = std::chrono::system_clock::now();
  const auto duration = current_time.time_since_epoch();
  const auto seconds =
      std::chrono::duration_cast<std::chrono::seconds>(duration).count();

  const ClickEvent click_event(23.4, 12.34, seconds);
  const InvokeCallback invoke_callback(123, click_event);

  NanoPack::Writer writer;
  const size_t bytes_written = invoke_callback.write_to(writer, 0);
  uint8_t *invoke_callback_serialized = writer.data();

  std::cout << "Raw bytes of ClickEvent: ";
  for (size_t i = 0; i < bytes_written; ++i) {
    std::cout << +invoke_callback_serialized[i] << " ";
  }

  std::cout << std::endl;
  std::cout << "Total bytes: " << bytes_written << std::endl;
  std::cout << "=========================" << std::endl;

  NanoPack::Reader reader(invoke_callback_serialized);

  int bytes_read;
  InvokeCallback invoke_callback_parsed;
  invoke_callback_parsed.read_from(reader);
  std::cout << "callback handle: " << invoke_callback_parsed.handle
            << std::endl;

  ClickEvent click_event_parsed;
  auto arg_reader = invoke_callback_parsed.args.into_reader();
  click_event_parsed.read_from(arg_reader);
  std::cout << "click event x, y: " << click_event_parsed.x << ", "
            << click_event_parsed.y << std::endl;
  std::cout << "timestamp: " << click_event_parsed.timestamp << std::endl;
}
