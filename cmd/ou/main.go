package main

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	log.SetFlags(0)
	app := cli.NewApp()
	app.Name = "ou"
	app.Usage = "orderup client"
	app.Version = "1.0.0"

	app.Commands = []cli.Command{
		{
			Name:   "target",
			Usage:  "save target server location",
			Action: target,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "port",
					Value: 5000,
					Usage: "orderup port",
				},
				cli.StringFlag{
					Name:  "passcode",
					Value: "",
					Usage: "orderup password",
				},
			},
		},
		{
			Name:   "list",
			Usage:  "get list of the pending orders in the queue",
			Action: list,
		},
		{
			Name:   "history",
			Usage:  "get list of the all finished orders in the queue",
			Action: history,
		},
	}

	app.Run(os.Args)
}
