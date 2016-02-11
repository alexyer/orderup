package main

// All bot command handlers located here.

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/boltdb/bolt"
)

// create-restaurant command.
// create-restaurant [restaurant name]
func (o *Orderup) createRestaurant(cmd *Cmd) (string, bool, *OrderupError) {
	switch {
	case len(cmd.Args) == 0:
		return "", true, NewOrderupError("Restaurant name is not given.", ARG_ERR)
	case len(cmd.Args) != 1:
		return "", true, NewOrderupError("Spaces are not allowed in restaurant name.", ARG_ERR)
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
		return "", true, NewOrderupError(err.Error(), CMD_ERR)
	}

	return fmt.Sprintf("%s restaurant created.", name), true, nil
}

// delete-restaurant command.
// delete-restaurant [restaurant name]
func (o *Orderup) deleteRestaurant(cmd *Cmd) (string, bool, *OrderupError) {
	switch {
	case len(cmd.Args) == 0:
		return "", true, NewOrderupError("Restaurant name is not given.", ARG_ERR)
	case len(cmd.Args) != 1:
		return "", true, NewOrderupError("Spaces are not allowed in restaurant name.", ARG_ERR)
	}

	name := cmd.Args[0]

	err := o.db.Update(func(tx *bolt.Tx) (err error) {
		// Get bucket with restaurants.
		b := tx.Bucket([]byte(RESTAURANTS))

		// Create new bucket for the new restaurant.
		err = b.DeleteBucket([]byte(name))
		if err != nil {
			return errors.New("Restaurant does not exist.")
		}

		return err
	})

	if err != nil {
		return "", true, NewOrderupError(err.Error(), CMD_ERR)
	}

	return fmt.Sprintf("Restaurant: %s deleted.", name), true, nil
}

// create-order command.
// create-order [restaurant name] [@username] [order]
func (o *Orderup) createOrder(cmd *Cmd) (string, bool, *OrderupError) {
	var (
		username       string
		restaurantName string
		order          string
		id             uint64
		orderCount     int
	)
	switch {
	case len(cmd.Args) < 3:
		return "", true, NewOrderupError("Wrong arguments", ARG_ERR)
	}

	username = cmd.Args[1]
	if username[0] != '@' {
		return "", true, NewOrderupError("Missing username", ARG_ERR)
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
		return "", true, NewOrderupError(err.Error(), CMD_ERR)
	}

	return fmt.Sprintf("%s order %d for %s %s - order %s. There are %d orders ahead of you.",
		restaurantName, int(id), username, order, order, orderCount), true, nil
}

// list command
// list [restaurant name]
func (o *Orderup) list(cmd *Cmd) (string, bool, *OrderupError) {
	if len(cmd.Args) != 1 {
		return "", true, NewOrderupError("Wrong arguments", ARG_ERR)
	}

	restaurantName := cmd.Args[0]

	ordersList, err := o.getPendingOrderList([]byte(restaurantName))
	if err != nil {
		return "", true, NewOrderupError(err.Error(), CMD_ERR)
	}

	// Format orders list properly
	result := fmt.Sprintf("%s: history:\n", restaurantName)
	for _, order := range *ordersList {
		result += order.String() + "\n"
	}

	return result, true, nil
}

// history command
// history [restaurant name]
func (o *Orderup) history(cmd *Cmd) (string, bool, *OrderupError) {
	if len(cmd.Args) != 1 {
		return "", true, NewOrderupError("Wrong arguments", ARG_ERR)
	}

	restaurantName := cmd.Args[0]

	history, err := o.getHistoryList([]byte(restaurantName))
	if err != nil {
		return "", true, NewOrderupError(err.Error(), CMD_ERR)
	}

	// Format orders list properly
	result := fmt.Sprintf("%s: history:\n", restaurantName)
	for _, order := range *history {
		result += order.String() + "\n"
	}

	return result, true, nil
}

// finish-order command
// finish-order [restaurant name] [order id]
func (o *Orderup) finishOrder(cmd *Cmd) (string, bool, *OrderupError) {
	var (
		order     Order
		orderData []byte
	)

	if len(cmd.Args) != 2 {
		return "", true, NewOrderupError("Wrong arguments", ARG_ERR)
	}

	restaurantName := cmd.Args[0]
	orderId, err := strconv.Atoi(cmd.Args[1])
	if err != nil {
		return "", true, NewOrderupError(err.Error(), CMD_ERR)
	}

	err = o.db.Batch(func(tx *bolt.Tx) (err error) {
		// Get bucket with restaurants.
		b := tx.Bucket([]byte(RESTAURANTS))

		restaurant := b.Bucket([]byte(restaurantName))
		if restaurant == nil {
			return errors.New(fmt.Sprintf("Restaurant %s does not exist", restaurantName))
		}

		orders := restaurant.Bucket([]byte(ORDERLIST))
		history := restaurant.Bucket([]byte(HISTORY))

		if orderId > orders.Stats().KeyN+1 {
			return errors.New("Too big order id. Order does not exist yet.")
		}

		orderData = orders.Get(itob(orderId))
		if orderData == nil {
			return errors.New("Order is already finished.")
		}

		// Delete order from the orders list
		if err := orders.Delete(itob(orderId)); err != nil {
			return err
		}

		// Put order in the history list
		return history.Put(itob(orderId), orderData)
	})

	if err != nil {
		return "", true, NewOrderupError(err.Error(), CMD_ERR)
	}

	if err := json.Unmarshal(orderData, &order); err != nil {
		return "", true, NewOrderupError(err.Error(), CMD_ERR)
	}

	return fmt.Sprintf("%s your order is finished. %s: Order: %d. %s",
		order.Username, restaurantName, order.Id, order.Order), true, nil
}

// help command.
func (o *Orderup) help(cmd *Cmd) (string, bool, *OrderupError) {
	return `Available commands:
				/orderup create-restaurant [name] -- Create a list of order numbers for restaurant name.
				/orderup delete-restaurant [name] -- Delete restaurant name and all orders in that restaurant.
				/orderup create-order [restaurant name] [@username] [order] -- Create a new order.
				/orderup finish-order [restaurant name]  [order id] -- Finish order.
				/orderup history [restaurant name] -- Show history for restaurant name.
				/orderup list [restaurant name] -- Get the list of orders for restaurant name.`, true, nil
}

// Helper function. Return error message with help contents.
func (o *Orderup) errorMessage(msg string) string {
	help, _, _ := o.help(nil)
	return fmt.Sprintf("%s\n%s", msg, help)
}

// Convert int to 8-byte big endian representation.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
