// AUTOMATICALLY GENERATED BY NANOPACK. DO NOT MODIFY BY HAND.

#ifndef MONTH_ENUM_NP_HXX
#define MONTH_ENUM_NP_HXX

#include <array>
#include <stdexcept>

class Month {
public:
  enum MonthMember {
    JANUARY,
    FEBRUARY,
    MARCH,
    APRIL,
    MAY,
    JUNE,
    JULY,
    AUGUST,
    SEPTEMBER,
    OCTOBER,
    NOVEMBER,
    DECEMBER,
  };

private:
  constexpr static std::array<int8_t, 12> values = {
      1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12,
  };
  MonthMember enum_value;
  int8_t _value;

public:
  Month() = default;

  explicit Month(const int8_t &value) {
    switch (value) {
    case 1:
      enum_value = JANUARY;
      break;
    case 2:
      enum_value = FEBRUARY;
      break;
    case 3:
      enum_value = MARCH;
      break;
    case 4:
      enum_value = APRIL;
      break;
    case 5:
      enum_value = MAY;
      break;
    case 6:
      enum_value = JUNE;
      break;
    case 7:
      enum_value = JULY;
      break;
    case 8:
      enum_value = AUGUST;
      break;
    case 9:
      enum_value = SEPTEMBER;
      break;
    case 10:
      enum_value = OCTOBER;
      break;
    case 11:
      enum_value = NOVEMBER;
      break;
    case 12:
      enum_value = DECEMBER;
      break;
    default:
      throw std::runtime_error("invalid value for enum Month");
    }
    _value = values[enum_value];
  }

  constexpr Month(MonthMember member)
      : enum_value(member), _value(values[member]) {}

  [[nodiscard]] constexpr const int8_t &value() const { return _value; }

  constexpr operator MonthMember() const { return enum_value; }

  explicit operator bool() const = delete;
};

#endif