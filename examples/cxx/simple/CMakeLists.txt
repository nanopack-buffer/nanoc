cmake_minimum_required(VERSION 3.25.2)

project(nanoc_simple_example
        VERSION 1.0.0
        LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED true)
set(CMAKE_EXPORT_COMPILE_COMMANDS true)

get_filename_component(SHARED_LIB_DIR ../lib ABSOLUTE)

add_executable(simple
		nanopack_message_factory.np.cxx
		simple_message.np.cxx
		nested_message.np.cxx
        main.cxx)

add_subdirectory(${SHARED_LIB_DIR}/cxx-nanopack ./cxx-nanopack)

target_include_directories(simple PUBLIC ${SHARED_LIB_DIR}/nanopack/include)
target_link_libraries(simple nanopack)
