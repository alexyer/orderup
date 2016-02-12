package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

var FINISH_ORDER_API_ENDPOINT = "/queues/orders/finish"
var FINISH_ORDER_API_METHOD = "PUT"

func finishOrder(c *cli.Context) {
	if len(c.Args()) != 2 {
		fmt.Println(WrongArgsError(c))
		return
	}

	name := c.Args()[0]
	id := c.Args()[1]

	resp, err := APICall(FINISH_ORDER_API_ENDPOINT, FINISH_ORDER_API_METHOD, struct {
		Name string `json:"name"`
		Id   string `json:"id"`
	}{
		Name: name,
		Id:   id,
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%s your order is finished. %s: Order %d. %s.\n",
		resp.Order.Username, name, resp.Order.Id, resp.Order.Order)
}
