# EFIBootSelector

This app allows you to change the default EFI boot order and to set your next EFI boot. Must run with admin privileges in order for the app to change EFI variables.

![](https://github.com/TaiPhamD/EFIBootSelector/blob/main/efi_boot.png)

## usage

- Run app and allow admin priviledge. You may figure out how to auto start this app on your own.
- This app has 2 basic functions:
   - Restart To: This will change EFI BootNext variable and will immediately restart the computer. It will not change your default boot.
   - Change Default Boot: This will change EFI BootOrder variable only. The check box on the sub menu entry indicates your current default boot. This selection will not cause a restart of the computer.

## future work

- Build client that's compatible with Linux 
- Build client that's compatible with MacOS. I am not sure if this is possible since I have not seen any macOS api that allows manipulation of EFI BootNext,BootOrder variables.
 (Would love to hear feedback if there's away)
- Persist elevated privilege in subsequent runs
- options to auto startup 

## build from source

## Prerequisite

- VS2019
- Golang 1.6+
- CMake v3+

run build.bat and it should create binaries in build/dist folder
