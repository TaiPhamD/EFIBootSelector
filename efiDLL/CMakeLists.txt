cmake_minimum_required(VERSION 3.0.0)
add_definitions(-DUNICODE -D_UNICODE)
FILE(GLOB MyHeaders ${CMAKE_CURRENT_SOURCE_DIR}/*.h)
FILE(GLOB MySources ${CMAKE_CURRENT_SOURCE_DIR}/*.cpp)
# Create a target for the library
ADD_LIBRARY(efiDLL SHARED ${MySources} ${MyHeaders})


#pragma comment(lib, "user32.lib")
#pragma comment(lib, "advapi32.lib")


target_link_libraries(efiDLL user32.lib  advapi32.lib)
set_target_properties(efiDLL 
                    PROPERTIES 
                    LINKER_LANGUAGE CXX
                    LIBRARY_OUTPUT_DIRECTORY "${CMAKE_BINARY_DIR}/dist"
                    RUNTIME_OUTPUT_DIRECTORY "${CMAKE_BINARY_DIR}/dist")


