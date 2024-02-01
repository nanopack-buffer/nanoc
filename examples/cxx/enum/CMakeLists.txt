cmake_minimum_required(VERSION 3.25.2)

project(nanoc_enum_example
        VERSION 1.0.0
        LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED True)

get_filename_component(SHARED_LIB_DIR ../lib ABSOLUTE)

add_executable(enum
        main.cxx
        column.np.cxx
        np_date.np.cxx)

add_subdirectory(${SHARED_LIB_DIR}/nanopack ./nanopack)

target_include_directories(enum PUBLIC ${SHARED_LIB_DIR}/nanopack/include)
target_link_libraries(enum nanopack)
