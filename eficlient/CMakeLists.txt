
project(efiBootSelector VERSION 0.1.0)

set(TARGET efiClient)
#set(GOPATH ${CMAKE_CURRENT_BINARY_DIR})
set(SRCS client.go)

# Custom go install 3rd party "Service" lib and build command              
add_custom_command(OUTPUT golang_eficlient_build.out
        DEPENDS ${SRCS}
        WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
        #COMMAND ${CMAKE_COMMAND} -E echo "Installing Golang 3rd party lib"
        #COMMAND go get github.com/kardianos/service
        #COMMAND go get github.com/getlantern/systray
        #COMMAND go get github.com/taiphamd/efibootselector
        COMMAND go mod download
        COMMAND ${CMAKE_COMMAND} -E echo "Compiling Golang efiClient project"
        #COMMAND set GOOS="windows"
        #COMMAND set GOARCH="386"
        COMMAND go build -ldflags "-H=windowsgui" -o "${CMAKE_BINARY_DIR}/dist"
        # -ldflags "-H=windowsgui"
)   


add_custom_target(efiClient ALL
        ${CMAKE_COMMAND} -E echo "Completed GoLang Compile"
        DEPENDS golang_eficlient_build.out)


