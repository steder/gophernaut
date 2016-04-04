package main

import "fmt"
import "net/http"
import "net/http/httputil"
import "net/url"
import "os/exec"
import "strings"

var hostname = fmt.Sprintf("http://127.0.0.1:%d", 8080)
var executable string = fmt.Sprintf("python -m SimpleHTTPServer %d", 8080)

func start_process(events chan int) {
	command_parts := strings.Split(executable, " ")
	command := exec.Command(command_parts[0], command_parts[1:]...)
	fmt.Printf("Command: %v\n", command)
	command.Start()
	_, ok := <-events
	if !ok {
		command.Process.Kill()
	}
}

func main() {
	events_channel := make(chan int)
	go start_process(events_channel)
	fmt.Printf("Gophernaut is ready for eBusiness!\n")
	hostname_url, _ := url.Parse(hostname)
	proxy := httputil.NewSingleHostReverseProxy(hostname_url)
	// Time to listen for connections!
	http.ListenAndServe(":8483", proxy)

}
