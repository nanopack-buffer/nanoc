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
    std::cout << "Total bytes: " << invoke_callback_data.size() << std::endl << std::endl;
}
