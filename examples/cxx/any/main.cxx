#include <iostream>
#include <vector>
#include <cstdint>
#include <chrono>

#include "click_event.np.hxx"
#include "invoke_callback.np.hxx"

int main()
{
    std::cout << "The any keyword in NanoPack allows you to store any NanoPack data type." << std::endl;

    const auto current_time = std::chrono::system_clock::now();
    const auto duration = current_time.time_since_epoch();
    const auto seconds = std::chrono::duration_cast<std::chrono::seconds>(duration).count();

    const ClickEvent click_event(23.4, 12.34, seconds);
    const InvokeCallback invoke_callback(123, click_event);
    const std::vector<uint8_t> invoke_callback_data = invoke_callback.data();

    std::cout << "Raw bytes of ClickEvent: ";
    for (const uint8_t b : invoke_callback_data)
    {
        std::cout << +b << " ";
    }
    std::cout << std::endl;
    std::cout << "Total bytes: " << invoke_callback_data.size() << std::endl;
    std::cout << "=========================" << std::endl;

    int bytes_read;
    const InvokeCallback invoke_callback_parsed(invoke_callback_data.begin(), bytes_read);
    std::cout << "callback handle: " << invoke_callback_parsed.handle << std::endl;

    const ClickEvent click_event_parsed(invoke_callback_parsed.args.as_reader(), bytes_read);
    std::cout << "click event x, y: " << click_event_parsed.x << ", " << click_event_parsed.y << std::endl;
    std::cout << "timestamp: " << click_event_parsed.timestamp << std::endl;
}
