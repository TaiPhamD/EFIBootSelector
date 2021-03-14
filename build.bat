@SETLOCAL
@REM alternative command to generate for vs2019:cd build; cmake .. -G "Visual Studio 16 2019" -A x64
@REM delete output folder 
if exist build rd /s /q build
mkdir build
cd build


SET PATH=C:\msys64\mingw64\bin;C:\msys64\usr\bin;%PATH%
SET CMAKE_C_COMPILER=gcc
SET CMAKE_CXX_COMPILER=g++
cmake -G "Unix Makefiles" -DCMAKE_BUILD_TYPE=Release ..
@REM cmake -DCMAKE_BUILD_TYPE=Release -G "NMake Makefiles" ..
@IF ERRORLEVEL 1 EXIT /B 1
cmake --build . --config Release
xcopy ..\install.bat dist\
xcopy ..\uninstall.bat dist\