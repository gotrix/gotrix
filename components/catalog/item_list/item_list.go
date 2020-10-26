package main

import (
	"fmt"

	"github.com/gotrix/gotrix"
)

var Component = helloWorld{}

type helloWorld struct{}

func (plugin *helloWorld) Include(cnf gotrix.ComponentParams) (string, error) {
	return fmt.Sprintf("%s component", cnf.Name()), nil
}
