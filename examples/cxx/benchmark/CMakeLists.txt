cmake_minimum_required(VERSION 3.25.2)

project(nanopack_benchmark
        VERSION 1.0.0
        LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED True)

get_filename_component(SHARED_LIB_DIR ../lib ABSOLUTE)

include(FindProtobuf)
find_package(Protobuf REQUIRED)
include_directories(${PROTOBUF_INCLUDE_DIR})

add_executable(benchmark
    main.cxx
    person.np.cxx
		person.pb.cc)

add_subdirectory(${SHARED_LIB_DIR}/cxx-nanopack ./cxx-nanopack)

target_include_directories(benchmark PUBLIC ${SHARED_LIB_DIR}/nanopack/include)
target_link_libraries(benchmark nanopack ${PROTOBUF_LIBRARY})
