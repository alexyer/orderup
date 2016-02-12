package main

import (
	"encoding/json"
	"io/ioutil"
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

// Persist credentils to disk.
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

// Read credentials from disk.
func ReadCredentials() (*Credentials, error) {
	f, err := os.Open(CRED_FILE)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	// Read file contents to the buffer.
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	cred := Credentials{}

	if err := json.Unmarshal(buf, &cred); err != nil {
		return nil, err
	}

	return &cred, nil
}
