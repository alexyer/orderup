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
			Usage:  "[queue name] -- Get the list of orders for queue name",
			Action: list,
		},
		{
			Name:   "history",
			Usage:  "[queue name] -- Show history for queue name",
			Action: history,
		},
		{
			Name:   "create-q",
			Usage:  "[name] -- Create a list of order numbers for queue <name>",
			Action: createQueue,
		},
		{
			Name:   "delete-q",
			Usage:  "[name] -- Delete queue <name> and all orders in that queue",
			Action: deleteQueue,
		},
		{
			Name:   "finish-order",
			Usage:  "[queue name]  [order id] -- Finish order",
			Action: finishOrder,
		},
		{
			Name:   "create-order",
			Usage:  "[queue name] [@username] [order] -- Create a new order",
			Action: createOrder,
		},
	}

	app.Run(os.Args)
}
