#include <nanopack/message.hxx>
#include <iostream>
#include <memory>
#include <cstdint>

#include "make_widget.np.hxx"
#include "text.np.hxx"
#include "widget.np.hxx"

int main() {
    std::cout << "NanoPack supports inheritance!" << std::endl;
    std::cout << "In this example, Text inherits Widget, so it inherits all the declared fields of Widgets." << std::endl;

    const std::unique_ptr<Widget> widget = std::make_unique<Text>(123, "hello world");
    std::vector<uint8_t> bytes = widget->data();

    int bytes_read;
    const std::unique_ptr<Widget> widget_from_factory = make_widget(bytes.begin(), bytes_read);

    std::cout << "Read " << bytes_read << " bytes" << std::endl;
    std::cout << "ID of text (inherited from Widget): " << widget_from_factory->id << std::endl;

    const Text *text = dynamic_cast<Text *>(widget_from_factory.get());
    std::cout << "Content of text (declared): " << text->content << std::endl;
}
