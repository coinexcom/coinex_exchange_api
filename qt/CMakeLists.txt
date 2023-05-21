cmake_minimum_required(VERSION 3.5)

project(coinex_bot LANGUAGES CXX)

set(CMAKE_INCLUDE_CURRENT_DIR ON)

set(CMAKE_AUTOUIC ON)
set(CMAKE_AUTOMOC ON)
set(CMAKE_AUTORCC ON)

set(CMAKE_CXX_STANDARD 11)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

find_package(Qt5Core)
find_package(Qt5WebSockets)

add_executable(coinex_bot
    main.cpp
    coinexapi.h coinexapi.cpp
    )
target_link_libraries(coinex_bot Qt5::Core Qt5::WebSockets)
