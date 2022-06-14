package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

func main() {
	cmd := exec.Command(os.Args[1])
	cmd.Args = os.Args[1:]
	stderr, err := cmd.StderrPipe()
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println("stderr fails")
	}
	if err := cmd.Start(); err != nil {
		log.Fatal("gnoty info:", err)
	}
	results, _ := io.ReadAll(stdout)
	fmt.Printf("%s", results)
	errors, _ := io.ReadAll(stderr)
	fmt.Printf("%s", errors)
	if err := cmd.Wait(); err != nil {
		log.Fatal(err.Error())
	}
}
