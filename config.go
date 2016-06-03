package gophernaut

import (
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
