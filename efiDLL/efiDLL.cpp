// shutdownDLL.cpp : Defines the exported functions for the DLL application.
// code adopted from: https://serverfault.com/questions/813695/how-do-i-stop-windows-10-install-from-modifying-bios-boot-settings

#include <windows.h>
#include <iomanip>
#include <iostream>
#include <sstream>
#include <locale>
#include <codecvt>

#ifdef _MSC_VER
#pragma comment(lib, "user32.lib")
#pragma comment(lib, "advapi32.lib")
#endif


  // Set Global UEFI GUID
const TCHAR globalGUID[] = TEXT("{8BE4DF61-93CA-11D2-AA0D-00E098032B8C}");
const TCHAR BootOrderStr[10] = TEXT("BootOrder");
const TCHAR BootNextStr[10] = TEXT("BootNext");

struct CloseHandleHelper
{
  void operator()(void *p) const { CloseHandle(p); }
};

/** Function to obtain required priviledges to issue shutdown or restart**/
BOOL SetPrivilege(HANDLE process, LPCWSTR name, BOOL on)
{
  HANDLE token;
  if (!OpenProcessToken(process, TOKEN_ADJUST_PRIVILEGES, &token))
    return FALSE;
  std::unique_ptr<void, CloseHandleHelper> tokenLifetime(token);
  TOKEN_PRIVILEGES tp;
  tp.PrivilegeCount = 1;
  if (!LookupPrivilegeValueW(NULL, name, &tp.Privileges[0].Luid))
    return FALSE;
  tp.Privileges[0].Attributes = on ? SE_PRIVILEGE_ENABLED : 0;
  return AdjustTokenPrivileges(token, FALSE, &tp, sizeof(tp), NULL, NULL);
}

/**Shutdown function**/
void shutdown(uint16_t *mode)
{
  //MODE 1 - shutdown 0 - restart
  if (*mode == 1)
    InitiateSystemShutdownEx(NULL, NULL, 0, FALSE, FALSE, 0);
  else
    InitiateSystemShutdownEx(NULL, NULL, 2, FALSE, TRUE, 0);
}

void changeBoot(uint16_t *data, uint16_t *mode)
{
  // MODE 0 : only change BootNext ( temporary next boot change)
  // MODE 1 : change BootOrder (permanent EFI boot order change)
  // Update UEFI
  const int bootOrderBytes = 2;
  const TCHAR(*bootOrderName)[10];

  if (*mode == 0)
    bootOrderName = &BootNextStr;
  else
    bootOrderName = &BootOrderStr;

  DWORD bootOrderAttributes = 7; // VARIABLE_ATTRIBUTE_NON_VOLATILE |
                                 // VARIABLE_ATTRIBUTE_BOOTSERVICE_ACCESS |
                                 // VARIABLE_ATTRIBUTE_RUNTIME_ACCESS
  SetFirmwareEnvironmentVariableEx(*bootOrderName, 
                                   globalGUID,
                                   data,
                                   bootOrderBytes,
                                   bootOrderAttributes);
}

extern "C"
{
  
  __declspec(dllexport) bool SystemShutdown(uint16_t *mode)
  {
    SetPrivilege(GetCurrentProcess(), SE_SHUTDOWN_NAME, TRUE);
    // we are just doign a normal shutdown
    shutdown(mode);
    // shutdown was successful
    return true;
  }

  __declspec(dllexport) bool SystemGetPermission()
  {
    // Get required shutdown priviledges
    SetPrivilege(GetCurrentProcess(), SE_SHUTDOWN_NAME, TRUE);
    SetPrivilege(GetCurrentProcess(), SE_SYSTEM_ENVIRONMENT_NAME, TRUE);
    return true;
  }

  __declspec(dllexport) bool SystemGetCurrentBoot(uint16_t *data, uint16_t *mode) {
    //Mode 0: BootNext, 1: BootOrder
    SetPrivilege(GetCurrentProcess(), SE_SYSTEM_ENVIRONMENT_NAME, TRUE);
    const TCHAR(*bootOrderName)[10];
    if (*mode == 0)
      bootOrderName = &BootNextStr;
    else
      bootOrderName = &BootOrderStr;
    const int BUFFER_SIZE = 4096;
    BYTE bootOrderBuffer[BUFFER_SIZE];
    DWORD bootOrderLength = 0;
    DWORD bootOrderAttributes = 7; // VARIABLE_ATTRIBUTE_NON_VOLATILE |
                                   // VARIABLE_ATTRIBUTE_BOOTSERVICE_ACCESS |
                                   // VARIABLE_ATTRIBUTE_RUNTIME_ACCESS
    bootOrderLength = GetFirmwareEnvironmentVariableEx(
        *bootOrderName, globalGUID, bootOrderBuffer, BUFFER_SIZE,
        &bootOrderAttributes);
    if (bootOrderLength == 0) {
      return false;
    }
    
    memcpy((char *)data, bootOrderBuffer, 2);
    return true;
  }

  __declspec(dllexport) bool GetBootEntries(char *output,uint16_t *size) {
    // Get boot entries
    SetPrivilege(GetCurrentProcess(), SE_SYSTEM_ENVIRONMENT_NAME, TRUE);
    std::string entries;
    std::stringstream my_ss;
    //setup wstring to string converter
    using convert_type = std::codecvt_utf8<wchar_t>;
    std::wstring_convert<convert_type, wchar_t> converter;
    DWORD bootOrderLength = 0;
    const int BUFFER_SIZE = 4096;

    for (DWORD i = 0; i < 10000; i++) {
      std::wstringstream bootOptionNameBuilder;
      bootOptionNameBuilder << "Boot" << std::uppercase << std::setfill(L'0')
                            << std::setw(4) << std::hex << (uint16_t)i;

      std::wstring bootOptionName(bootOptionNameBuilder.str());
      BYTE bootOptionInfoBuffer[BUFFER_SIZE];
      DWORD bootOptionInfoLength = GetFirmwareEnvironmentVariableEx(
          bootOptionName.c_str(), globalGUID, bootOptionInfoBuffer, BUFFER_SIZE,
          nullptr);
      if (bootOptionInfoLength == 0) {
        /*std::cout << "Failed getting option info for option at offset " << i
                  << std::endl;*/
        continue;
      }
      uint32_t *bootOptionInfoAttributes =
          reinterpret_cast<uint32_t *>(bootOptionInfoBuffer);

      std::wstring description(reinterpret_cast<wchar_t *>(bootOptionInfoBuffer + sizeof(uint32_t) + sizeof(uint16_t)));
      /*
      std::wcout << "Boot option name:" << bootOptionName << std::endl;
      std::wcout << "Boot description: " << description
                  << " integer ID: " << std::dec << i << std::endl
                  << std::endl;
      */
       my_ss <<  std::dec << i << ":" <<  converter.to_bytes(description) << std::endl;
    }

    entries = my_ss.str();
    size_t max_size = std::min(entries.size()+1,static_cast<std::string::size_type>(*size));
    memcpy(output,entries.c_str(),max_size);

    //entries.c_str();
    return true;
  }

  __declspec(dllexport) bool SystemChangeBoot(uint16_t *data, uint16_t *mode)
  {
    SetPrivilege(GetCurrentProcess(), SE_SYSTEM_ENVIRONMENT_NAME, TRUE);
    // data is boot integer id
    // Mode 0: BootNext, 1: BootOrder
    changeBoot(data,mode);
    return true;
  }
}
