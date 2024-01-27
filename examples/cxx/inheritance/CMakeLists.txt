cmake_minimum_required(VERSION 3.25.2)

project(nanoc_inheritance_example
        VERSION 1.0.0
        LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED True)

get_filename_component(SHARED_LIB_DIR ../lib ABSOLUTE)

add_executable(inheritance
        main.cxx
        widget.np.cxx
        make_widget.np.cxx
        text.np.cxx
        nanopack_message_factory.np.cxx)

add_subdirectory(${SHARED_LIB_DIR}/nanopack ./nanopack)

target_include_directories(inheritance PUBLIC ${SHARED_LIB_DIR}/nanopack/include)
target_link_libraries(inheritance nanopack)
