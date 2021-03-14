# EFIBootSelector

This app allows you to change the default EFI boot order and to set your next EFI boot.

![](https://github.com/TaiPhamD/EFIBootSelector/blob/main/EFIBOOTSELECTOR.jpg)

## usage

### installation 
- Run install.bat. Installation will do the following:
  - Install client & server app to C:\efibootselector
  - Create new windows service called EFIBootSelectorService 
  - Add registry key to auto start client app : ADD HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Run /v EFIBootSelector

- This app has 2 basic functions:
   - Restart To: This will change EFI BootNext variable and will immediately restart the computer. It will not change your default boot.
   - Change Default Boot: This will change EFI BootOrder variable only. The check box on the sub menu entry indicates your current default boot. This selection will not cause a restart of the computer.

## future work

- Build client that's compatible with Linux 
- Build client that's compatible with MacOS. I am not sure if this is possible since I have not seen any macOS api that allows manipulation of EFI BootNext,BootOrder variables.
 (Would love to hear feedback if there's away)

## build from source

## Prerequisite
It will most likely work with other versions and just need to adopt build.bat to change generator for your build system. Below are the specific versions that I tested the build script with:
- MSYS2 for GCC tool chain (https://www.msys2.org/)
   - install the following packages:
      - pacman -S --needed base-devel mingw-w64-x86_64-toolchain 
      - pacman -S cmake msys2-w32api-headers msys2-w32api-runtime
- Golang 1.16+
- CMake v3+

run build.bat and it should create binaries in build/dist folder
