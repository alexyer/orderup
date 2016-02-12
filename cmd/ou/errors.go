package main

import (
	"errors"
	"fmt"

	"github.com/codegangsta/cli"
)

func WrongArgsError(c *cli.Context) error {
	return errors.New(fmt.Sprintf("Wrong arguments.\n%s %s", c.Command.Name, c.Command.Usage))
}
