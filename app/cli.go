package app

import (
	"os"
	"os/exec"
	"strconv"

	"github.com/mkideal/cli"
)

type StartT struct {
	Port int `cli:"port" dft:"8080" usage:"Application port to listen"`
}

func StartFromCLI(ctx *cli.Context) error {
	argv := ctx.Argv().(*StartT)
	cmd := exec.Command("go",
		"run", "cmd/app/main.go",
		"--port", strconv.Itoa(argv.Port))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
