package main

import (
	"fmt"
	"os"

	"github.com/TarsCloud/TarsGo/tars"
	"github.com/TarsCloud/TarsGo/tars/util/rogger"

	"github.com/TarsDemo/Tars-MiniProgramm-Service-ClubActivityServer/tars-protocol/LifeService"
)

var comm *tars.Communicator

//SLOG 日志
var SLOG = rogger.GetLogger("ServerLog")

func main() {
	comm = tars.NewCommunicator()
	// Get server config
	cfg := tars.GetServerConfig()

	// New servant imp
	imp := new(ClubActivityManagerImp)
	err := imp.Init()
	if err != nil {
		fmt.Printf("ClubActivityManagerImp init fail, err:(%s)\n", err)
		os.Exit(-1)
	}
	// New servant
	app := new(LifeService.ClubActivityManager)
	// Register Servant
	app.AddServantWithContext(imp, cfg.App+"."+cfg.Server+".ClubActivityManagerObj")

	// Run application
	tars.Run()
}
