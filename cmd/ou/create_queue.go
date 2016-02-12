package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

var CREATE_Q_API_ENDPOINT = "/queues"
var CREATE_Q_API_METHOD = "POST"

func createQueue(c *cli.Context) {
	if len(c.Args()) != 1 {
		fmt.Println(WrongArgsError(c))
		return
	}

	name := c.Args()[0]

	resp, err := APICall(CREATE_Q_API_ENDPOINT, CREATE_Q_API_METHOD, struct {
		Name string `json:"name"`
	}{
		Name: name,
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(resp.Response)
}
