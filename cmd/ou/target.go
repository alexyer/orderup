package main

// target command.

import (
	"fmt"
	"log"
	"net/http"

	"github.com/codegangsta/cli"
)

func target(c *cli.Context) {
	if len(c.Args()) < 1 {
		log.Fatal("error: hostname required")
	}

	host := c.Args()[0]
	port := c.Int("port")
	passcode := c.String("passcode")

	// Try to connect to the host
	client := http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/api/v1/queues/orders/list", host, port), nil)
	req.SetBasicAuth("", passcode)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	// Check if passcode is correct
	if resp.StatusCode == http.StatusUnauthorized {
		log.Fatal("Wrong passcode. Check passcode and try again.")
	}

	err = WriteCredentials(&Credentials{
		Host:     host,
		Port:     port,
		Passcode: passcode,
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Print("Authorization successful.")
}
