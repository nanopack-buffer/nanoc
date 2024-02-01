#ifndef TEST_ENUM_HXX
#define TEST_ENUM_HXX

#include <string>

class Test
{
public:
    enum Member : uint8_t
    {
        Member1,
        Member2,
        Member3
    };

private:
    const std::array<std::string, 3> values = {{"test", "test2", "test3"}};
    inline const static Test lookup[3] = {Member1, Member2, Member3};

    Member enum_value;
    std::string _value;

public:
    static Test from_value(const uint8_t value)
    {
        return lookup[value];
    };

    const std::string& value = _value;

    Test() = default;

    constexpr Test(Member member) : enum_value(member), value(values[member])
    {
    }

    constexpr operator Member() const { return enum_value; }

    // Prevent usage: if(fruit)
    explicit operator bool() const = delete;
};

#endif //TEST_ENUM_HXX
