package app

import (
	"github.com/gotrix/gotrix"
	"github.com/mkideal/cli"
)

type StartT struct {
	Port int `cli:"port" dft:"8080" usage:"Application port to listen"`
}

func StartFromCLI(ctx *cli.Context) error {
	argv := ctx.Argv().(*StartT)
	srv, err := New(&gotrix.AppConfig{
		Port: argv.Port,
	})
	if err != nil {
		return err
	}
	return srv.Run()
}
