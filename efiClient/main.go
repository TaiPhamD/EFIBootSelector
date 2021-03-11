package main

import "C"
import (
	"eficlient/icondata"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	"github.com/getlantern/systray"
	"golang.org/x/sys/windows"
)

var GetPermissionFunc *syscall.Proc
var ShutDownFunc *syscall.Proc
var ChangeBootFunc *syscall.Proc
var GetBootEntriesFunc *syscall.Proc
var GetCurrentBootFunc *syscall.Proc
var BootArray []BootEntry

var mode uint16
var data uint16

type BootEntry struct {
	index    uint16
	id       string
	selected bool
	nb       *systray.MenuItem //next boot object
	db       *systray.MenuItem //default boot object
}

func runMeElevated() {
	verb := "runas"
	exe, _ := os.Executable()
	cwd, _ := os.Getwd()
	args := strings.Join(os.Args[1:], " ")

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	argPtr, _ := syscall.UTF16PtrFromString(args)

	var showCmd int32 = 1 //SW_NORMAL

	err := windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(0)
}

func amAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		fmt.Println("admin no")
		return false
	}
	fmt.Println("admin yes")
	return true
}

func main() {
	if !amAdmin() {
		runMeElevated()
	}
	systray.Run(onReady, onExit)
}

func getEntries() []BootEntry {

	var default_boot uint16
	mode = 1
	// Get current default boot
	GetCurrentBootFunc.Call(uintptr(unsafe.Pointer(&default_boot)), uintptr(unsafe.Pointer(&mode)))

	// Get all available boot entries
	var buffer_size uint16
	buffer_size = 2048
	c_string_buffer := make([]byte, buffer_size)
	GetBootEntriesFunc.Call(uintptr(unsafe.Pointer(&c_string_buffer)), uintptr(unsafe.Pointer(&buffer_size)))
	entries := C.GoString((*C.char)(unsafe.Pointer(&c_string_buffer)))
	//fmt.Printf("%s\n", entries)
	list := strings.Split(entries, "\n")
	result := make([]BootEntry, 0, len(list))
	for _, s := range list {
		pair := strings.Split(s, ":")
		if len(pair) < 2 {
			continue
		}

		value, _ := strconv.ParseUint(pair[0], 10, 32)
		var temp_bool bool
		if uint16(value) == default_boot {
			temp_bool = true
		} else {
			temp_bool = false
		}

		result = append(result, BootEntry{uint16(value), pair[1], temp_bool, nil, nil})
	}
	return result
}

func setDefault(myentry *[]BootEntry, index uint16) {
	//clear all existing entries
	for i := 0; i < len(*myentry); i++ {
		s := &(*myentry)[i]
		s.db.Uncheck()
		if s.index == index {
			s.db.Check()
		}
	}
}

func onReady() {

	// Load DLL function for windows api
	efifunc := syscall.MustLoadDLL("efiDLL")
	GetBootEntriesFunc = efifunc.MustFindProc("GetBootEntries")
	GetCurrentBootFunc = efifunc.MustFindProc("SystemGetCurrentBoot")
	// GetPermissionFunc.Call()
	BootArray = getEntries()

	// Build systray GUI
	systray.SetIcon(icondata.MainIcon)
	systray.SetTitle("EFI Boot Selector")
	systray.SetTooltip("EFI Boot Selector")
	mRestart := systray.AddMenuItem("Restart To", "Restart to next selected boot")
	mSetDefault := systray.AddMenuItem("Change Boot Default", "Change next default boot")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	//mQuit.SetIcon(icondata.MainIcon)
	go func() {
		<-mQuit.ClickedCh
		//fmt.Println("Requesting quit")
		systray.Quit()
		//fmt.Println("Finished quitting")
	}()

	//fmt.Println(mRestart)
	//for _, s := range BootArray {
	for i := 0; i < len(BootArray); i++ {
		//fmt.Println(i, s.index, s.id, s.selected)
		s := &BootArray[i]
		s.nb = mRestart.AddSubMenuItem(s.id, "Restart to new OS on next boot")
		s.db = mSetDefault.AddSubMenuItem(s.id, "Change default boot")
		// Add check box if it's a default boot entry
		if s.selected {
			s.db.Check()
		}
		go func() {
			for {
				select {
				case <-s.db.ClickedCh:
					fmt.Println("We clicked default boot to:", s.index)
					setDefault(&BootArray, s.index)
					mode = 1 //for default boot change
					data = s.index
					loaddll := syscall.MustLoadDLL("efiDLL")
					//defer loaddll.Release()
					ChangeBootFunc := loaddll.MustFindProc("SystemChangeBoot")
					ChangeBootFunc.Call(uintptr(unsafe.Pointer(&data)), uintptr(unsafe.Pointer(&mode)))
				case <-s.nb.ClickedCh:
					fmt.Println("We are restarting to")
					mode = 0 //for default boot change
					data = s.index
					loaddll := syscall.MustLoadDLL("efiDLL")
					//defer loaddll.Release()
					ChangeBootFunc := loaddll.MustFindProc("SystemChangeBoot")
					RestartFunc := loaddll.MustFindProc("SystemShutdown")
					ChangeBootFunc.Call(uintptr(unsafe.Pointer(&data)), uintptr(unsafe.Pointer(&mode)))
					RestartFunc.Call(uintptr(unsafe.Pointer(&mode)))
				}
			}
		}()

	}

	fmt.Println(mSetDefault)
	// Sets the icon of a menu item. Only available on Mac and Windows.

}

func onExit() {
	// clean up here
}
