cmake_minimum_required(VERSION 3.0)

#include_directories(/env/gopath/src/github.com/PaddlePaddle/Paddle/paddle/go/cclient/build/)

cc_test(cgo_test
        SRCS
        main.c
        DEPS
        cclient)

# add_executable(main main.c)
# add_dependencies(main cclient)
# set (CMAKE_EXE_LINKER_FLAGS "-pthread")
# target_link_libraries(main cclient)
