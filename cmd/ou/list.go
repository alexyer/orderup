package main

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"
)

var LIST_API_ENDPOINT = "/queues/orders/list"
var LIST_API_METHOD = "GET"

func list(c *cli.Context) {
	if len(c.Args()) != 1 {
		log.Fatal("Wrong arguments")
	}

	name := c.Args()[0]

	resp, err := APICall(LIST_API_ENDPOINT, LIST_API_METHOD, struct {
		Name string `json:"name"`
	}{
		Name: name,
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s: pending orders:\n", name)
	for _, order := range resp.Orders {
		fmt.Printf("\t%s\n", order)
	}
}
