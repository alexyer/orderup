package main

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"
)

var DELETE_Q_API_ENDPOINT = "/queues"
var DELETE_Q_API_METHOD = "DELETE"

func deleteQueue(c *cli.Context) {
	if len(c.Args()) != 1 {
		log.Fatal("Wrong arguments")
	}

	name := c.Args()[0]

	resp, err := APICall(DELETE_Q_API_ENDPOINT, DELETE_Q_API_METHOD, struct {
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
