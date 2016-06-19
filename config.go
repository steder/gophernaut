package gophernaut

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

/*
Config is a Gophernaught config structure used to parse gophernaut.conf
*/
type Config struct {
	Host string
	Port int
	/*
	    	Debug bool
	   	Pool struct {
	   		Size int
	   		Template struct {
	   			Name string
	   		}
	   	}
	*/
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

/*
ReadConfig reads gophernaut.conf and returns a Config object
*/
func ReadConfig() *Config {
	data, error := ioutil.ReadFile("etc/template.conf")
	check(error)
	c := Config{}
	yaml.Unmarshal(data, &c)
	return &c
}

// GetExecutables uses our config to provide a set of processes to start
func (c *Config) GetExecutables() []string {
	var executables = []string{
		fmt.Sprintf("python -m SimpleHTTPServer %d", 8080),
		fmt.Sprintf("python -m SimpleHTTPServer %d", 8081),
	}
	return executables
}

// GetHostnames uses our config to provide a set of hostnames to dispatch requests to
func (c *Config) GetHostnames() []string {
	var hostnames = []string{
		fmt.Sprintf("http://127.0.0.1:%d", 8080),
		fmt.Sprintf("http://127.0.0.1:%d", 8081),
	}
	return hostnames
}
