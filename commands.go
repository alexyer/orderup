package main

// All bot command handlers located here.

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/boltdb/bolt"
)

// create-restaurant command.
// create-restaurant [restaurant name]
func (o *Orderup) createRestaurant(cmd *Cmd) string {
	switch {
	case len(cmd.Args) == 0:
		return "Restaurant name is not given."
	case len(cmd.Args) != 1:
		return "Spaces are not allowed in restaurant name."
	}

	name := cmd.Args[0]

	err := o.db.Update(func(tx *bolt.Tx) (err error) {
		// Get bucket with restaurants.
		b := tx.Bucket([]byte(RESTAURANTS))

		// Create new bucket for the new restaurant.
		r, err := b.CreateBucket([]byte(name))
		if err != nil {
			return errors.New("Restaurant already exists.")
		}

		// Create 2 subbucktes.
		// One for pending orders, another for finished orders.
		_, err = r.CreateBucket([]byte(ORDERLIST))
		_, err = r.CreateBucket([]byte(HISTORY))

		return err
	})

	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("%s restaurant created.", name)
}

// create-order command.
// create-order [restaurant name] [@username] [order]
func (o *Orderup) createOrder(cmd *Cmd) string {
	var (
		username       string
		restaurantName string
		order          string
		id             uint64
		orderCount     int
	)
	switch {
	case len(cmd.Args) < 3:
		return o.errorMessage("Wrong arguments")
	}

	username = cmd.Args[1]
	if username[0] != '@' {
		return o.errorMessage("Missing username")
	}

	restaurantName = cmd.Args[0]

	err := o.db.Update(func(tx *bolt.Tx) (err error) {
		// Get bucket with restaurants.
		b := tx.Bucket([]byte(RESTAURANTS))

		restaurant := b.Bucket([]byte(restaurantName))
		if restaurant == nil {
			return errors.New(fmt.Sprintf("Restaurant %s does not exist", restaurantName))
		}

		orders := restaurant.Bucket([]byte(ORDERLIST))

		// Prepare order data
		id, _ = orders.NextSequence()
		order = strings.Join(cmd.Args[2:], " ")
		orderCount = orders.Stats().KeyN

		// JSON serialize order
		buf, err := json.Marshal(&Order{
			Username: username,
			Order:    order,
			Id:       int(id),
		})

		if err != nil {
			return err
		}

		// Store order into the database
		return orders.Put(itob(int(id)), buf)
	})

	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("%s order %d for %s %s - order %s. There are %d orders ahead of you.",
		restaurantName, int(id), username, order, order, orderCount)
}

// list command
// list [restaurant name]
func (o *Orderup) list(cmd *Cmd) string {
	var ordersList []Order

	if len(cmd.Args) != 1 {
		return o.errorMessage("Wrong arguments")
	}

	restaurantName := cmd.Args[0]

	err := o.db.View(func(tx *bolt.Tx) (err error) {
		// Get bucket with restaurants.
		b := tx.Bucket([]byte(RESTAURANTS))

		restaurant := b.Bucket([]byte(restaurantName))
		if restaurant == nil {
			return errors.New(fmt.Sprintf("Restaurant %s does not exist", restaurantName))
		}

		orders := restaurant.Bucket([]byte(ORDERLIST))
		c := orders.Cursor()

		// Iterate over all orders, decode and store in the orders list
		for k, v := c.First(); k != nil; k, v = c.Next() {
			order := Order{}
			if err := json.Unmarshal(v, &order); err != nil {
				return err
			}

			ordersList = append(ordersList, order)
		}

		return
	})

	if err != nil {
		return err.Error()
	}

	// Format orders list properly
	result := ""
	for _, order := range ordersList {
		result += order.String() + "\n"
	}

	return result
}

// help command.
func (o *Orderup) help(cmd *Cmd) string {
	return `Available commands:
				/orderup create-restaurant [name] -- Create a list of order numbers for restaurant name.
				/orderup create-order [restaurant name] [@username] [order] -- Create a new order.
				/orderup list [restaurant name] -- Get the list of orders for restaurant name.`
}

// Helper function. Return error message with help contents.
func (o *Orderup) errorMessage(msg string) string {
	return fmt.Sprintf("%s\n%s", msg, o.help(nil))
}

// Convert int to 8-byte big endian representation.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
