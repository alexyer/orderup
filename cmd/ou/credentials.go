package main

import (
	"encoding/json"
	"log"
	"os"
	"os/user"
)

var CRED_FILE string // Credentials file location

func init() {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	CRED_FILE = user.HomeDir + "/.orderup"
}

type Credentials struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Passcode string `json:"passcode"`
}

func WriteCredentials(c *Credentials) error {
	f, err := os.Create(CRED_FILE)
	if err != nil {
		return err
	}

	defer f.Close()

	buf, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	_, err = f.Write(buf)

	return err
}
