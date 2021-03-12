package model

// common struct for RPC server and client
type Args struct {
	Data  uint16
	Mode  uint16
	Token string
}

// On startup the server will randomize token and port
// and save this data to current path of server app
// the client app will need to use this data to connect to RPC
type ServerInfo struct {
	Port  uint16
	Token string
}
