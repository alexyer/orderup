package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	host     string
	port     int
	db       string
	password string
)

func init() {
	flag.StringVar(&host, "host", "localhost", "orderup host")
	flag.IntVar(&port, "port", 5000, "orderup port")
	flag.StringVar(&password, "passcode", "", "password")
}

// Check connection
func checkConn() {
	client := http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/api/v1/queues/orders/list", host, port), nil)
	req.SetBasicAuth("", password)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		log.Fatal("Wrong passcode. Check passcode and try again.")
	}
}

func main() {
	flag.Parse()
	checkConn()
}
