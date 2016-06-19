package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/steder/gophernaut"
)

func main() {
	log.SetPrefix("gophernaut ")
	log.SetFlags(log.Ldate | log.Ltime)
	c := gophernaut.ReadConfig()
	log.Printf("Host %s and Port %d\n", c.Host, c.Port)
	log.Printf("Pool %d, %s, %s, %s\n",
		c.Pool.Size,
		c.Pool.Template.Name,
		c.Pool.Template.Executable,
		c.Pool.Template.Hostname,
	)

	pool := gophernaut.Pool{Executables: c.GetExecutables()}
	pool.Start()
	go pool.ManageProcesses()

	log.Printf("Gophernaut is gopher launch!\n")
	// TODO: our own ReverseProxy implementation of at least, ServeHTTP so that we can
	// monitor the response codes to track successes and failures
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", c.Port),
		http.HandlerFunc(gophernaut.GetGopherHandler(c.GetHostnames()))))
}
