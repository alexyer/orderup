package main

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"
)

var HISTORY_API_ENDPOINT = "/queues/orders/history"
var HISTORY_API_METHOD = "GET"

func history(c *cli.Context) {
	if len(c.Args()) != 1 {
		log.Fatal("Wrong arguments")
	}

	name := c.Args()[0]

	resp, err := APICall(HISTORY_API_ENDPOINT, HISTORY_API_METHOD, struct {
		Name string `json:"name"`
	}{
		Name: name,
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%s: history:\n", name)
	for _, order := range resp.Orders {
		fmt.Printf("\t%s\n", order)
	}
}
