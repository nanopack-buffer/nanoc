#include <iostream>

#include "np_date.np.hxx"
#include "column.np.hxx"

int main()
{
    std::cout << "A simple program demonstrating enums in NanoPack." << std::endl;
    std::cout << "In NanoPack, enums are serialized into their backing types." << std::endl;
    std::cout <<
        "The backing type can be specified in the schema. When none is specified, nanoc will determine the most appropriate type."
        << std::endl << std::endl;

    const NpDate date(19, Week::MONDAY, Month::JUNE, 2000);
    const std::vector<uint8_t> date_data = date.data();
    std::cout << "The Date message uses the Week and Month enum, both of which are backed by an int8." << std::endl;
    std::cout << "Raw bytes of Date:";
    for (const uint8_t b : date_data)
    {
        std::cout << +b << " ";
    }
    std::cout << std::endl;
    std::cout << "Total bytes: " << date_data.size() << std::endl << std::endl;

    const Column column(Alignment::CENTER);
    const std::vector<uint8_t> column_data = column.data();
    std::cout << "The Column message uses the Alignment enum which is backed by a string." << std::endl;
    std::cout << "Raw bytes of Column:";
    for (const uint8_t b : column_data)
    {
        std::cout << +b << " ";
    }
    std::cout << std::endl;
    std::cout << "Total bytes: " << column_data.size() << std::endl << std::endl;
}
