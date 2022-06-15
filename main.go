package main

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	bolt "go.etcd.io/bbolt"
	"io"
	"os"
	"os/exec"
	"pigil/internal/database"
	"pigil/internal/types"
	"pigil/internal/utils"
	"time"
)

const mainTag = "main"

func main() {

	db, err := bolt.Open(utils.DatabaseName, 0666, nil)
	if err != nil {
		utils.ErrorLogger(err, mainTag)
	}

	err = godotenv.Load()
	if err != nil {
		utils.ErrorLogger(errors.New("cannot load .env file"), mainTag)
	}
	boltLocalDb := database.NewBoltDbService(db, utils.LocalBucket)
	boltConfigDb := database.NewBoltDbService(db, utils.ConfigBucket)
	dbService := database.NewDatabaseService(boltLocalDb, boltConfigDb)

	if os.Args[1] == "bumf" {
		cliHandler(os.Args, dbService)
	} else {
		executor(os.Args, dbService)
	}
	defer db.Close()

}

func cliHandler(args []string, service database.Service) {
	switch args[2] {
	case utils.CliDb:
		ListCommand(service)
	case utils.CliAuth:
		GoogleAuth(service)
	case utils.CliStatus:
		Status(service)
	case utils.CliLogout:
		Logout(service)
	}
}

func executor(args []string, service database.Service) {
	cmd := exec.Command(args[1])
	cmd.Args = args[1:]
	stderr, err := cmd.StderrPipe()
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		utils.ErrorInformation("stderr pipe broke")
	}
	start := time.Now()
	if err := cmd.Start(); err != nil {
		utils.ErrorLogger(err, mainTag)
	}
	ci := types.CommandInformation{
		CommandName:      args[1],
		CommandArguments: args[2:],
		TimeOfExecution:  time.Now().UnixMicro(),
		WasSuccessful:    true,
	}
	results, _ := io.ReadAll(stdout)
	fmt.Printf("%s", results)
	errorOutput, _ := io.ReadAll(stderr)
	fmt.Printf("%s", errorOutput)
	if err := cmd.Wait(); err != nil {
		ci.WasSuccessful = false
		Notify(service)
		utils.ErrorInformation(err.Error())
	}
	end := time.Now()
	life := end.Sub(start)
	utils.GreenPrinter(fmt.Sprintf("runtime: %f seconds", life.Seconds()))
	ci.ExecutionTime = life.Seconds()
	InsertCommand(service, ci)
}
