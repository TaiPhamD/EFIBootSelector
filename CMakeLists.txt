cmake_minimum_required(VERSION 3.0.0)
project(efiBootSelector VERSION 0.1.0)

set(CPACK_PROJECT_NAME ${PROJECT_NAME})
set(CPACK_PROJECT_VERSION ${PROJECT_VERSION})
include(CPack)

add_subdirectory("efidll")
add_subdirectory("eficlient")
add_subdirectory("efiserver")
