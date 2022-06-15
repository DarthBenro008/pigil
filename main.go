package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"gnoty/internal/database"
	"gnoty/internal/types"
	"gnoty/internal/utils"
	bolt "go.etcd.io/bbolt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {

	db, err := bolt.Open(utils.DatabaseName, 0666, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
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
		log.Println("stderr fails")
	}
	if err := cmd.Start(); err != nil {
		log.Fatal("gnoty info:", err)
	}
	ci := types.CommandInformation{
		CommandName:      args[1],
		CommandArguments: args[2:],
		ExecutionTime:    time.Now().UnixMicro(),
		WasSuccessful:    true,
	}
	results, _ := io.ReadAll(stdout)
	fmt.Printf("%s", results)
	errors, _ := io.ReadAll(stderr)
	fmt.Printf("%s", errors)
	if err := cmd.Wait(); err != nil {
		ci.WasSuccessful = false
		//log.Fatal(err.Error())
	}
	InsertCommand(service, ci)
}
