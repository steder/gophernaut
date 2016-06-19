package main

import (
	"log"
	"net/http"

	"github.com/steder/gophernaut"
)

func main() {
	log.SetPrefix("gophernaut ")
	log.SetFlags(log.Ldate | log.Ltime)
	c := gophernaut.ReadConfig()
	log.Printf("Host %s and Port %d\n", c.Host, c.Port)

	pool := gophernaut.Pool{Executables: c.GetExecutables()}
	pool.Start()
	go pool.ManageProcesses()

	log.Printf("Gophernaut is gopher launch!\n")
	// TODO: our own ReverseProxy implementation of at least, ServeHTTP so that we can
	// monitor the response codes to track successes and failures
	log.Fatal(http.ListenAndServe(":8483",
		http.HandlerFunc(gophernaut.GetGopherHandler(c.GetHostnames()))))
}
