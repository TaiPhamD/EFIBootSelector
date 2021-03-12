package api

import "C"
import (
	"errors"
	"fmt"
	"syscall"
	"unsafe"

	"github.com/getlantern/systray"
	model "github.com/taiphamd/efibootselector/common"
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

type EFI_RPC int

var MyServerInfo model.ServerInfo 
func Authenticate(args *model.Args) bool {
	return (*args).Token == MyServerInfo.Token
}

func (t *EFI_RPC) GetEntries(args *model.Args, result *string) error {
	if !Authenticate(args){
		return errors.New("invalid token")
	}
	efifunc := syscall.MustLoadDLL("efiDLL")
	GetBootEntriesFunc = efifunc.MustFindProc("GetBootEntries")
	var buffer_size uint16
	buffer_size = 2048
	c_string_buffer := make([]byte, buffer_size)
	GetBootEntriesFunc.Call(uintptr(unsafe.Pointer(&c_string_buffer)), uintptr(unsafe.Pointer(&buffer_size)))
	*result = C.GoString((*C.char)(unsafe.Pointer(&c_string_buffer)))
	//fmt.Println("RPC GetEntries")
	//fmt.Println(*result)
	return nil
}

func (t *EFI_RPC) GetDefaultBoot(args *model.Args, result *uint16) error {
	if !Authenticate(args){
		return errors.New("invalid token")
	}	
	// Get current default boot
	var default_boot uint16
	mode = 1
	efifunc := syscall.MustLoadDLL("efiDLL")
	GetCurrentBootFunc = efifunc.MustFindProc("SystemGetCurrentBoot")
	GetCurrentBootFunc.Call(uintptr(unsafe.Pointer(&default_boot)), uintptr(unsafe.Pointer(&mode)))
	*result = default_boot
	//fmt.Println("RPC GetDefaultBootCalled")
	//fmt.Println(default_boot)
	return nil
}

func (t *EFI_RPC) SetDefaultBoot(args *model.Args, result *uint16) error {
	if !Authenticate(args){
		return errors.New("invalid token")
	}
	fmt.Println((*args).Token)
	// Get current default boot
	mode = (*args).Mode
	data = (*args).Data
	loaddll := syscall.MustLoadDLL("efiDLL")
	//defer loaddll.Release()
	ChangeBootFunc := loaddll.MustFindProc("SystemChangeBoot")
	ChangeBootFunc.Call(uintptr(unsafe.Pointer(&data)), uintptr(unsafe.Pointer(&mode)))
	return nil
}

func (t *EFI_RPC) ShutDown(args *model.Args, result *uint16) error {
	if !Authenticate(args){
		return errors.New("invalid token")
	}
	mode = (*args).Mode
	efifunc := syscall.MustLoadDLL("efiDLL")
	ShutdownFunc := efifunc.MustFindProc("SystemShutdown")
	ShutdownFunc.Call(uintptr(unsafe.Pointer(&mode)))
	return nil
}
