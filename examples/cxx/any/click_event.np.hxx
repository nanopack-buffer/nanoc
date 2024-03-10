// AUTOMATICALLY GENERATED BY NANOPACK. DO NOT MODIFY BY HAND.

#ifndef CLICK_EVENT_NP_HXX
#define CLICK_EVENT_NP_HXX

#include <nanopack/message.hxx>
#include <nanopack/nanopack.hxx>
#include <nanopack/reader.hxx>
#include <vector>

struct ClickEvent : NanoPack::Message {
  static constexpr NanoPack::TypeId TYPE_ID = 837166865;

  double x;
  double y;
  int64_t timestamp;

  ClickEvent() = default;

  ClickEvent(double x, double y, int64_t timestamp);

  ClickEvent(std::vector<uint8_t>::const_iterator begin, int &bytes_read);

  ClickEvent(const NanoPack::Reader &reader, int &bytes_read);

  size_t write_to(std::vector<uint8_t> &buf, int offset) const override;

  [[nodiscard]] NanoPack::TypeId type_id() const override;

  [[nodiscard]] int header_size() const override;

  [[nodiscard]] std::vector<uint8_t> data() const override;
};

#endif
