cmake_minimum_required(VERSION 3.0)

get_filename_component(PARENT_DIR ${CMAKE_CURRENT_SOURCE_DIR} DIRECTORY)
get_filename_component(PARENT_DIR ${PARENT_DIR} DIRECTORY)
set(CMAKE_MODULE_PATH ${CMAKE_MODULE_PATH} "${PARENT_DIR}/cmake")

project(cxx_go C Go)

include(golang)
include(flags)

go_library(paddle_master SHARED)

if(PROJ_ROOT)
  # using missing file paddle_master_copy as output, so this command
  # will run everytime.
  add_custom_command(OUTPUT paddle_master_copy
    COMMAND cp ${CMAKE_CURRENT_BINARY_DIR}/${paddle_master_LIB_NAME} ${PROJ_ROOT}/python/paddle/v2/master/
    DEPENDS paddle_master)
  add_custom_target(paddle_master_python ALL DEPENDS paddle_master_copy)
endif(PROJ_ROOT)
