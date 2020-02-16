package main

import (
	"github.com/urfave/cli"
	"testing"
)

func TestStartServer(t *testing.T) {
	err := StartServer(&cli.Context{
		Context: nil,
	})
	if err != nil {
		t.Errorf(err.Error())
	}

}
func TestPrintHello(t *testing.T) {
	PrintHello()
}
