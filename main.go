package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
)

var (
	host string
	port int
	db   string
)

func init() {
	flag.StringVar(&host, "host", "localhost", "host to listen")
	flag.IntVar(&port, "port", 5000, "port to listen")
	flag.StringVar(&db, "db", "orderup.db", "database file")
}

func main() {
	flag.Parse()

	// Create new instance of bot.
	bot, err := NewOrderup(db)
	if err != nil {
		log.Fatal(err)
	}

	// Configure shutdown process.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for _ = range c {
			bot.Shutdown()
			os.Exit(0)
		}
	}()

	log.Printf("Orderup started on %s:%d.\n", host, port)

	http.HandleFunc("/orderup", bot.RequestHandler)
	http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil)
}
