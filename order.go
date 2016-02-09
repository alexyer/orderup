package main

import "fmt"

type Order struct {
	Username string `json:username`
	Order    string `json:order`
	Id       int    `json:id`
}

func (o Order) String() string {
	return fmt.Sprintf("%d %s - %s", o.Id, o.Username, o.Order)
}
