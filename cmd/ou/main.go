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
	}

	app.Run(os.Args)
}
