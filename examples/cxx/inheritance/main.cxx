#include <cstddef>
#include <cstdint>
#include <iostream>
#include <memory>
#include <nanopack/message.hxx>

#include "make_widget.np.hxx"
#include "nanopack/reader.hxx"
#include "nanopack/writer.hxx"
#include "text.np.hxx"
#include "widget.np.hxx"

int main() {
  std::cout << "NanoPack supports inheritance!" << std::endl;
  std::cout << "In this example, Text inherits Widget, so it inherits all the "
               "declared fields of Widgets."
            << std::endl;

  std::shared_ptr<Text> text = std::make_shared<Text>(123, "hello world");

  std::cout << "test" << std::endl;

  // IMPORTANT!!
  // don't forget to call write_to on the correct type, otherwise not all data
  // will be serialized! for example, if write_to were called on widget instead
  // of text, Widget::write_to would have been called instead of Text::write_to
  // which means fields declared in Text would not have been serialized!!
  NanoPack::Writer writer;
  const size_t bytes_written = text->write_to(writer, 0);
  uint8_t *buf = writer.data();

  // now, we will deserialize from the buffer using a Reader
  // make_widget is an automatically generated function that can create the
  // correct instance of a message based on the type id specified in the raw
  // buffer.

  size_t bytes_read;
  NanoPack::Reader reader(buf);
  std::unique_ptr<Widget> widget_from_factory = make_widget(reader, bytes_read);

  std::cout << "Read " << bytes_read << " bytes" << std::endl;
  std::cout << "ID of text (inherited from Widget): " << widget_from_factory->id
            << std::endl;

  // to access the fields of Text, we need to cast the message to a Text

  const Text *text_deserialized =
      dynamic_cast<Text *>(widget_from_factory.get());
  std::cout << "Content of text: " << text_deserialized->content << std::endl;
}
