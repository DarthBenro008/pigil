package main

import (
	_ "embed"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/DarthBenro008/pigil/internal/database"
	"github.com/DarthBenro008/pigil/internal/types"
	"github.com/DarthBenro008/pigil/internal/utils"
	bolt "go.etcd.io/bbolt"
)

const mainTag = "main"

var Version = "development"

//go:embed secrets.txt
var secrets string

func main() {
	dirname, err := os.UserHomeDir()
	if err != nil {
		utils.ErrorLogger(err, mainTag)
	}
	secretsArray := strings.Split(secrets, " ")
	utils.GoogleClientId = secretsArray[0]
	utils.GoogleClientSecret = secretsArray[1]
	dirname = fmt.Sprintf("%s/%s", dirname, ".pigil")
	if _, err := os.Stat(dirname); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(dirname, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
	dirname = fmt.Sprintf("%s/%s", dirname, utils.DatabaseName)
	db, err := bolt.Open(dirname, 0666, nil)
	if err != nil {
		utils.ErrorLogger(err, mainTag)
	}

	boltLocalDb := database.NewBoltDbService(db, utils.LocalBucket)
	boltConfigDb := database.NewBoltDbService(db, utils.ConfigBucket)
	dbService := database.NewDatabaseService(boltLocalDb, boltConfigDb)
	IsFirstTime(dbService)

	if len(os.Args) == 1 {
		Help()
		os.Exit(0)
	}

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
	case utils.CliDiscord:
		if len(args) == 2 {
			DiscordToggle(service, true)
		} else if args[3] == "disable" {
			DiscordToggle(service, true)
		} else {
			utils.InformationLogger(
				"To disable discord webhook run `pigil bumf discord disable`")
		}

	case utils.CliHelp:
		Help()
	default:
		utils.InformationLogger("Invalid pigil Command!")
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
	err = cmd.Wait()
	end := time.Now()
	life := end.Sub(start)
	utils.GreenPrinter(fmt.Sprintf("runtime: %f seconds", life.Seconds()))
	ci.ExecutionTime = life.Seconds()
	if err != nil {
		ci.WasSuccessful = false
		Notify(service, ci)
		utils.ErrorInformation(err.Error())
	}

	InsertCommand(service, ci)
}
