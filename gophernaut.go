package main

import "os/exec"
import "fmt"
import "strings"

var executable string = "python -m SimpleHTTPServer 8080"

func start_process() {
	command_parts := strings.Split(executable, " ")
	command := exec.Command(command_parts[0], command_parts[1:]...)
	fmt.Printf("Command: %v\n", command)
	command.Start()
}

func main() {
	go start_process()
	fmt.Printf("Gophernaut is ready for eBusiness!\n")
	for {
	}
}
