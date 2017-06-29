cmake_minimum_required(VERSION 3.0)

go_library(paddle_master SHARED)


# using missing file paddle_master_copy as output, so this command
# will run everytime.
add_custom_command(OUTPUT paddle_master_copy
  COMMAND cp ${CMAKE_CURRENT_BINARY_DIR}/${paddle_master_LIB_NAME} ${PROJ_ROOT}/python/paddle/v2/master/
  DEPENDS paddle_master)
add_custom_target(paddle_master_python ALL DEPENDS paddle_master_copy)
