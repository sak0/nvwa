package main

import (
	"os"
	"fmt"
	"github.com/urfave/cli"
)

var ( 
	usage = "Usage: ?"
	version = "0.0.1"
)

func main() {
	app := cli.NewApp()
	app.Name = "nvwa"
        app.Version = version
	app.Usage = usage

        app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "enable debug output for logging",
		},
	}

	app.Commands = []cli.Command{
		joinCommand,
	}

        app.Action = Run
	app.Run(os.Args)
}

func Run(ctx *cli.Context) {
	if ctx.Bool("debug") {
		fmt.Println("debug mode")
	}
	fmt.Println("Run.")
}
