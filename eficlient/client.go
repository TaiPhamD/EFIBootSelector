package main

import "C"
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/taiphamd/efibootselector/eficlient/icondata"

	"github.com/getlantern/systray"
	model "github.com/taiphamd/efibootselector/common"
)

var BootArray []BootEntry
var mode uint16
var data uint16
var rpc_client *rpc.Client
var serverinfo model.ServerInfo //use to validate authenticate with rpc server
type BootEntry struct {
	index    uint16
	id       string
	selected bool
	nb       *systray.MenuItem //next boot object
	db       *systray.MenuItem //default boot object
}

func main() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	//set up log file
	filelog, errlog := os.OpenFile(dir+"\\eficlient.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if errlog != nil {
		log.Fatal(errlog)
	}
	defer filelog.Close()
	log.SetOutput(filelog)

	//get server info
	jsonFile, err := os.Open(dir + "\\efiserver_config.json")
	if err != nil {
		log.Fatal(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &serverinfo)	


	systray.Run(onReady, onExit)
}

func getEntries() []BootEntry {

	rpc_client, err := rpc.DialHTTP("tcp", "localhost:"+strconv.Itoa(int(serverinfo.Port)))
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var args model.Args
	// Get default boot
	var default_boot uint16
	args.Mode = 1
	args.Token = serverinfo.Token
	err = rpc_client.Call("EFI_RPC.GetDefaultBoot", args, &default_boot)
	if err != nil {
		log.Fatal("EFI_RPC.GetEntries error:", err)
	}

	// Get boot entries
	args.Mode = 1
	var entries string
	err = rpc_client.Call("EFI_RPC.GetEntries", args, &entries)
	if err != nil {
		log.Fatal("EFI_RPC.GetEntries error:", err)
	}

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

	BootArray = getEntries()
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
					rpc_client, err := rpc.DialHTTP("tcp", "localhost:"+strconv.Itoa(int(serverinfo.Port)))
					if err != nil {
						log.Fatal("dialing:", err)
					}
					var args model.Args
					args.Mode = 1
					args.Data = s.index
					args.Token = serverinfo.Token
					err = rpc_client.Call("EFI_RPC.SetDefaultBoot", args, nil)
					if err != nil {
						log.Fatal("EFI_RPC.SetDefaultBoot error:", err)
					}

				case <-s.nb.ClickedCh:
					fmt.Println("We are restarting to")
					/*mode = 0 //for default boot change
					data = s.index*/
					rpc_client, err := rpc.DialHTTP("tcp", "localhost:"+strconv.Itoa(int(serverinfo.Port)))
					if err != nil {
						log.Fatal("dialing:", err)
					}
					var args model.Args
					args.Mode = 0
					args.Data = s.index
					args.Token = serverinfo.Token
					err = rpc_client.Call("EFI_RPC.SetDefaultBoot", args, nil)
					if err != nil {
						log.Fatal("EFI_RPC.SetDefaultBoot error:", err)
					}

					err = rpc_client.Call("EFI_RPC.ShutDown", args, nil)
					if err != nil {
						log.Fatal("EFI_RPC.ShutDown error:", err)
					}
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
