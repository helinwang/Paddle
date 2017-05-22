cmake_minimum_required(VERSION 3.0)

# find cmake directory modules
get_filename_component(PARENT_DIR ${CMAKE_CURRENT_SOURCE_DIR} DIRECTORY)
get_filename_component(PARENT_DIR ${PARENT_DIR} DIRECTORY)
get_filename_component(PARENT_DIR ${PARENT_DIR} DIRECTORY)
set(CMAKE_MODULE_PATH ${CMAKE_MODULE_PATH} "${PARENT_DIR}/cmake")

# if(GTEST_INCLUDE_DIR AND GTEST_LIBRARIES)
#   message("-- Found gtest (include: ${GTEST_INCLUDE_DIR}, library: ${GTEST_LIBRARIES})")
# else()	
  # find #include <majel/xx.h>
  get_filename_component(PARENT_DIR ${CMAKE_CURRENT_SOURCE_DIR} DIRECTORY)
  include_directories(${PARENT_DIR})

  # enable c++11
  set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -std=c++11")

  # enable gtest
  set(THIRD_PARTY_PATH ./third_party)
  set(WITH_TESTING ON)
  include(external/gtest)
#endif()

set(CMAKE_MODULE_PATH ${CMAKE_MODULE_PATH} "${CMAKE_CURRENT_SOURCE_DIR}/cmake")

project(cclient CXX C Go)

include(generic)
#include(util)

#go_extern(pserver github.com/PaddlePaddle/Paddle/paddle/go/pserver)
#go_library(cclient STATIC pserver)
go_library(cclient SRCS cclient.go)
#add_dependencies(cclient pserver)
#add_subdirectory(test)



get_cmake_property(_variableNames VARIABLES)
# foreach (_variableName ${_variableNames})
#     message(STATUS "${_variableName}=${${_variableName}}")
# endforeach()

message("lib:   " ${cclient_lib})

include_directories(${cclient_BINARY_DIR})
cc_test(cclient_test
        SRCS
        main.cc
        DEPS
        ${cclient_lib})
