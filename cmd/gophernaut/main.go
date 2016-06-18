package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"strings"

	"github.com/steder/gophernaut"
)

// Event is basically just an enum
type Event int

// Events that can be generated by our child processes
const (
	Start Event = iota
	Shutdown
	PiningForTheFjords
)

// TODO look into "go generate stringer -type Event"
func (e Event) String() string {
	return fmt.Sprintf("Event(%d)", e)
}

var hostname = fmt.Sprintf("http://127.0.0.1:%d", 8080)
var executable = fmt.Sprintf("python -m SimpleHTTPServer %d", 8080)

var hostnames = []string{
	fmt.Sprintf("http://127.0.0.1:%d", 8080),
	fmt.Sprintf("http://127.0.0.1:%d", 8081),
}

var executables = []string{
	fmt.Sprintf("python -m SimpleHTTPServer %d", 8080),
	fmt.Sprintf("python -m SimpleHTTPServer %d", 8081),
}

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

var requestCount = 0

func myHandler(w http.ResponseWriter, myReq *http.Request) {
	requestPath := myReq.URL.Path
	// TODO: multiprocess, pick one of n hostnames based on pool status
	hostname := hostnames[requestCount%2] // TODO get rid of this hard coded 2
	requestCount++
	targetURL, _ := url.Parse(hostname)
	director := func(req *http.Request) {
		targetQuery := targetURL.RawQuery
		req.URL.Scheme = targetURL.Scheme
		// TODO: adjust request host to assign the request to the appropriate child process
		req.URL.Host = targetURL.Host

		// clean up but preserve trailing slash:
		trailing := strings.HasSuffix(req.URL.Path, "/")
		req.URL.Path = path.Join(targetURL.Path, req.URL.Path)
		if trailing && !strings.HasSuffix(req.URL.Path, "/") {
			req.URL.Path += "/"
		}

		// preserve query string:
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
	}

	proxy := &httputil.ReverseProxy{Director: director}

	staticHandler := http.StripPrefix("/static", http.FileServer(http.Dir("static")))
	adminTemplate := template.Must(template.ParseFiles("templates/admin.html"))
	adminHandler := func(w http.ResponseWriter, req *http.Request) {
		adminTemplate.Execute(w, nil)
	}

	switch {
	case requestPath == "/admin":
		adminHandler(w, myReq)
		return
	case strings.HasPrefix(requestPath, "/static"):
		staticHandler.ServeHTTP(w, myReq)
		return
	}
	proxy.ServeHTTP(w, myReq)
}

func main() {
	log.SetPrefix("gophernaut ")
	log.SetFlags(log.Ldate | log.Ltime)
	c := gophernaut.ReadConfig()
	log.Printf("Host %s and Port %d\n", c.Host, c.Port)

	controlChannel := make(chan Event)
	eventsChannel := make(chan Event)

	// Handle signals to try to do a graceful shutdown:
	receivedSignals := make(chan os.Signal, 1)
	signal.Notify(receivedSignals, os.Interrupt) // , syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for sig := range receivedSignals {
			fmt.Printf("Received signal, %s, shutting down workers...\n", sig)
			break
		}
		close(controlChannel)
		signal.Stop(receivedSignals)
	}()

	// Actually start some processes
	for _, executable := range executables {
		go startProcess(controlChannel, eventsChannel, executable)
	}

	// wait for child processes to exit before shutting down:
	processCount := 2 // TODO get rid of these hard coded 2s!
	stoppedCount := 0
	go func() {
		for event := range eventsChannel {
			if event == Shutdown {
				stoppedCount++
			}
			if processCount == stoppedCount {
				fmt.Printf("%d workers stopped, shutting down.\n", processCount)
				os.Exit(1)
			}
		}
	}()

	log.Printf("Gophernaut is gopher launch!\n")
	// TODO: our own ReverseProxy implementation of at least, ServeHTTP so that we can
	// monitor the response codes to track successes and failures
	log.Fatal(http.ListenAndServe(":8483", http.HandlerFunc(myHandler)))
}
