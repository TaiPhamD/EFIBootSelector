@echo off
::  prompt user to get elevated priviledge for installation
::-------------------------------------
REM  --> Check for permissions
>nul 2>&1 "%SYSTEMROOT%\system32\cacls.exe" "%SYSTEMROOT%\system32\config\system"

REM --> If error flag set, we do not have admin.
if '%errorlevel%' NEQ '0' (
    echo Requesting administrative privileges...
    goto UACPrompt
) else ( goto gotAdmin )

:UACPrompt
    echo Set UAC = CreateObject^("Shell.Application"^) > "%temp%\getadmin.vbs"
    set params = %*:"="
    echo UAC.ShellExecute "cmd.exe", "/c %~s0 %params%", "", "runas", 1 >> "%temp%\getadmin.vbs"

    "%temp%\getadmin.vbs"
    del "%temp%\getadmin.vbs"
    exit /B

:gotAdmin
    pushd "%CD%"
    CD /D "%~dp0"
::--------------------------------------

::ENTER YOUR CODE BELOW:

:: Delete windows service service if it exists and create new one
SC QUERY EFIBootSelectorService > NUL
IF ERRORLEVEL 1060 GOTO MISSING
ECHO service exist stopping then deleting service
sc stop EFIBootSelectorService 
sc delete EFIBootSelectorService 
GOTO NOSERVICE

:MISSING
ECHO service does not exist creating new one...

:NOSERVICE

:: kill any existing client
QPROCESS "eficlient.exe">NUL 2> nul
IF %ERRORLEVEL% EQU 0 taskkill /IM eficlient.exe


timeout 2

::Copy binaries to target destination
if exist "c:\efibootselector\" rd /q /s "c:\efibootselector
mkdir c:\efibootselector
xcopy build\dist\efiserver.exe c:\efibootselector\
xcopy build\dist\eficlient.exe c:\efibootselector\
xcopy build\dist\efiDLL.dll c:\efibootselector\

::Create windows service
sc create EFIBootSelectorService binPath="C:\efibootselector\efiserver.exe"
sc config EFIBootSelectorService start= auto
sc start EFIBootSelectorService 

:: add registry key to auto start eficlient on windows startup
REG ADD HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Run /v EFIBootSelector /d "C:\efibootselector\eficlient.exe" /f

timeout 5
start C:\efibootselector\eficlient.exe
pause