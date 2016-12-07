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
	Pool struct {
		Size     int
		Template struct {
			Name       string
			Hostname   string
			Executable string
		}
	}
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
	var executables []string
	for x := 0; x < c.Pool.Size; x++ {
		executables = append(executables,
			fmt.Sprintf(c.Pool.Template.Executable, 8080+x))
	}
	return executables
}

// GetHostnames uses our config to provide a set of hostnames to dispatch requests to
func (c *Config) GetHostnames() []string {
	var hostnames []string
	for x := 0; x < c.Pool.Size; x++ {
		hostnames = append(hostnames,
			fmt.Sprintf(c.Pool.Template.Hostname, 8080+x))
	}
	return hostnames
}
