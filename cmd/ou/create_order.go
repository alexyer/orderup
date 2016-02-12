package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/codegangsta/cli"
)

var CREATE_ORDER_API_ENDPOINT = "/queues/order"
var CREATE_ORDER_API_METHOD = "POST"

func createOrder(c *cli.Context) {
	if len(c.Args()) < 3 {
		log.Fatal("Wrong arguments")
	}

	name := c.Args()[0]
	username := c.Args()[1]
	desc := strings.Join(c.Args()[2:], " ")

	if username[0] != '@' {
		fmt.Println("Missing username")
		return
	}

	resp, err := APICall(CREATE_ORDER_API_ENDPOINT, CREATE_ORDER_API_METHOD, struct {
		Name        string `json:"name"`
		User        string `json:"user"`
		Description string `json:"description"`
	}{
		Name:        name,
		User:        username[1:],
		Description: desc,
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%s order %d for %s - order %s. There are %d orders ahead of you.\n",
		name, resp.Order.Id, resp.Order.Username, resp.Order.Order, resp.OrdersAhead)
}
