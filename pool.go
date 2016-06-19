package gophernaut

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
)

// Event is basically just an enum
type Event int

// Events that can be generated by our child processes
//go:generate stringer -type=Event
const (
	Start Event = iota
	Shutdown
	PiningForTheFjords
)

func copyToLog(dst *log.Logger, src io.Reader) {
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		dst.Print(scanner.Text())
	}
}

func startProcess(control <-chan Event, events chan<- Event, executable string) {
	procLog := log.New(os.Stdout, fmt.Sprintf("gopher-worker(%s) ", executable), log.Ldate|log.Ltime)

	commandParts := strings.Split(executable, " ")
	command := exec.Command(commandParts[0], commandParts[1:]...)
	log.Printf("Command: %v\n", command)

	stdout, err := command.StdoutPipe()
	if err != nil {
		procLog.Fatalln("Unable to connect to stdout from command...")
	}
	stderr, err := command.StderrPipe()
	if err != nil {
		procLog.Fatalln("Unable to connect to stderr from command...")
	}

	//go io.Copy(os.Stdout, stdout)
	//go io.Copy(os.Stderr, stderr)
	go copyToLog(procLog, stdout)
	go copyToLog(procLog, stderr)
	command.Start()

	events <- Start
	for {
		_, ok := <-control
		if !ok {
			fmt.Println("Killing worker process after receiving close event.")
			command.Process.Kill()
			events <- Shutdown
			break
		}
	}
}

// Pool manages the pool of processes to which gophernaut dispatches
// requests.
type Pool struct {
	Executables []string

	stoppedCount   int
	processCount   int
	controlChannel chan Event
	eventsChannel  chan Event
}

// Start up the pool
func (p *Pool) Start() {
	p.controlChannel = make(chan Event)
	p.eventsChannel = make(chan Event)

	// Handle signals to try to do a graceful shutdown:
	receivedSignals := make(chan os.Signal, 1)
	signal.Notify(receivedSignals, os.Interrupt) // , syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for sig := range receivedSignals {
			fmt.Printf("Received signal, %s, shutting down workers...\n", sig)
			break
		}
		close(p.controlChannel)
		signal.Stop(receivedSignals)
	}()

	// Actually start some processes
	for _, executable := range p.Executables {
		go startProcess(p.controlChannel, p.eventsChannel, executable)
	}
}

// ManageProcesses manages soem processes
func (p *Pool) ManageProcesses() {
	for event := range p.eventsChannel {
		switch event {
		case Shutdown:
			p.stoppedCount++
		case Start:
			p.processCount++
		}
		if p.processCount == p.stoppedCount {
			log.Printf("%d workers stopped, shutting down.\n", p.processCount)
			os.Exit(1)
		}
	}
}
