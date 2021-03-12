@echo off
:: BatchGotAdmin
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

::Download tdm-gcc run time
ECHO Please install TDM-GCC runtime first: 
ECHO https://github.com/jmeubank/tdm-gcc/releases/download/v9.2.0-tdm64-1/tdm64-gcc-9.2.0.exe
@echo off
setlocal
:PROMPT
SET /P AREYOUSURE=Did you install required GCC runtime (Y/[N])?
IF /I "%AREYOUSURE%" NEQ "Y" GOTO END
SC QUERY EFIBootSelectorService > NUL
IF ERRORLEVEL 1060 GOTO MISSING
ECHO service exist stopping then deleting service
sc stop EFIBootSelectorService 
sc delete EFIBootSelectorService 
GOTO NOSERVICE

:MISSING
ECHO service does not exist creating new one...

:NOSERVICE

taskkill /IM eficlient.exe

timeout 2

if exist "c:\efibootselector\" rd /q /s "c:\efibootselector
mkdir c:\efibootselector
xcopy efiserver.exe c:\efibootselector\
xcopy eficlient.exe c:\efibootselector\
xcopy efiDLL.dll c:\efibootselector\


sc create EFIBootSelectorService binPath="C:\efibootselector\efiserver.exe"
sc config EFIBootSelectorService start= auto
sc start EFIBootSelectorService 

REG ADD HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Run /v EFIBootSelector /d "C:\Windows\System32\cmd.exe /k start C:\efibootselector\eficlient.exe" /f 

timeout 5
start C:\efibootselector\eficlient.exe

:END
endlocal

pause