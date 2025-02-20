cmake_minimum_required(VERSION 3.20)

project(haversine C)

set(CMAKE_C_STANDARD 11)

file(GLOB SOURCES "src/*.c")
file(GLOB HEADERS "src/*.h")

add_executable(haversine ${SOURCES} ${HEADERS})

target_compile_features(haversine PRIVATE c_std_11)
target_compile_options(haversine PRIVATE -Wall -Wextra -Werror)

target_link_libraries(haversine PRIVATE m)
