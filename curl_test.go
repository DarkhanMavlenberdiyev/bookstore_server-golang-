package main

import (
	"testing"
)

func TestPrintHello(t *testing.T) {
	res := PrintHello(2)
	if res!=4{
		t.Errorf("Not correct")
	}
}
