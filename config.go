package gophernaut

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

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

func ReadConfig() *Config {
	data, error := ioutil.ReadFile("etc/template.conf")
	check(error)
	c := Config{}
	yaml.Unmarshal(data, &c)
	return &c
}
