package gophernaut

import "testing"

func TestConfig(t *testing.T) {
	t.Log(Config{
		Host: "foo.bar",
		Port: 1234,
	})
}

func TestRead(t *testing.T) {
	c := ReadConfig()
	t.Log(c)
}
