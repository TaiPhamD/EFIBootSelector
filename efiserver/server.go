package main

import (
	"crypto/rand"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"path/filepath"

	"github.com/kardianos/service"
	model "github.com/taiphamd/efibootselector/common"
	"github.com/taiphamd/efibootselector/efiserver/api"
)

//type RPCHandler struct{}
type program struct{}

func main() {
	//Get file path from where the exe is launched
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	
	//set up log file
	filelog, errlog := os.OpenFile(dir+"\\efiserver.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 644)
	if errlog != nil {
		log.Fatal(errlog)
	}
	defer filelog.Close()
	log.SetOutput(filelog)

	svcConfig := &service.Config{
		Name:        "EFIBootSelector",
		DisplayName: "EFIBootSelector RPC Service",
		Description: "EFIBootSelector RPC Service",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}

}

func (p *program) Stop(s service.Service) error {
	log.Print("Stopped Shutdown service\n")
	// Stop should not block. Return with a few seconds.
	return nil
}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}

func (p *program) run() {

	efi_rpc := new(api.EFI_RPC)
	rpc.Register(efi_rpc)
	rpc.HandleHTTP()

	// Try to find any available port
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		log.Fatal("finding dynamic port", err)
		panic(err)
	}
	// Store it to %APPDATA% so client app can connect to it
	log.Print("Using port:", listener.Addr().(*net.TCPAddr).Port)
	var MyServerInfo model.ServerInfo 
	MyServerInfo.Port = uint16(listener.Addr().(*net.TCPAddr).Port)
	MyServerInfo.Token,_ = GenerateRandomString(50)

	api.MyServerInfo = MyServerInfo

	file, _ := json.MarshalIndent(MyServerInfo, "", " ")
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	err = ioutil.WriteFile(dir + "\\efiserver_config.json", file, 0644)
	if err != nil {
		log.Fatal("cannot save config json", err)
		panic(err)
	}	
	panic(http.Serve(listener, nil))
}

func GenerateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}