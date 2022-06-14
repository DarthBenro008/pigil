package main

import (
	"fmt"
	bolt "go.etcd.io/bbolt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {

	db, err := bolt.Open(databaseName, 0666, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	boltDatabase := NewBoldDbService(db)
	dbService := NewDatabaseService(boltDatabase)

	if os.Args[1] == "bumf" {
		cliHandler(os.Args, dbService)
	} else {
		executor(os.Args, dbService)
	}
	defer db.Close()

}

func cliHandler(args []string, service DatabaseService) {
	switch args[2] {
	case cliDb:
		ListCommand(service)
	}
}

func executor(args []string, service DatabaseService) {
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
	ci := CommandInformation{
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
