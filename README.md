# EFIBootSelector

This app allows you to change the default EFI boot order and to set your next EFI boot.

![](https://github.com/TaiPhamD/EFIBootSelector/blob/main/EFIBOOTSELECTOR.jpg)

## usage
### Prerequisite
 - Requires GCC runtime: https://github.com/jmeubank/tdm-gcc/releases/download/v9.2.0-tdm64-1/tdm64-gcc-9.2.0.exe
### installation 
- Run install.bat. See [release note](https://github.com/TaiPhamD/EFIBootSelector/releases/tag/v0.2.1) if you want to know what install and uninstall does.
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
- VS2019
- Golang 1.16+
- CMake v3+
- GCC Compiler something like:  https://github.com/jmeubank/tdm-gcc/releases/download/v9.2.0-tdm64-1/tdm64-gcc-9.2.0.exe

run build.bat and it should create binaries in build/dist folder
