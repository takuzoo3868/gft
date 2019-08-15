package main

import (
	"os"

	"github.com/urfave/cli"
	"github.com/takuzoo3868/gft/cmd"
)

func main() {
	app := cli.NewApp()
	app.Name = "gft"
	app.Usage = "file transfer with gRPC"
	app.Version = "0.1"
	app.Commands = []cli.Command{
		cmd.DownloadCommand(),
		cmd.ServeCommand(),
		cmd.ListCommand(),
	}
	app.Run(os.Args)
}