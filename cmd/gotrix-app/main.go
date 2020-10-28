package main

import (
	"flag"
	"log"

	"github.com/gotrix/gotrix"
	"github.com/gotrix/gotrix/app"
)

func main() {
	var (
		componentPaths gotrix.PathFlags
		templatePaths  gotrix.PathFlags
	)
	flag.Var(&componentPaths, "component-path", "path to component include dir")
	flag.Var(&templatePaths, "template-path", "path to template include dir")
	flag.Parse()
	log.Println(componentPaths)
	a, err := app.New(&gotrix.AppConfig{
		ComponentPaths: componentPaths,
	})
	if err != nil {
		log.Fatal(err)
	}
	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}
