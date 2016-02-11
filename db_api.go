package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/boltdb/bolt"
)

// DB calls.
// Available for Slack and REST APIs.

// Get orders list for the queue.
func (o *Orderup) getOrderList(queue []byte) (*[]Order, error) {
	var orderList []Order

	err := o.db.View(func(tx *bolt.Tx) (err error) {
		// Get bucket with restaurants.
		b := tx.Bucket([]byte(RESTAURANTS))

		restaurant := b.Bucket(queue)
		if restaurant == nil {
			return errors.New(fmt.Sprintf("Queue %s does not exist", queue))
		}

		orders := restaurant.Bucket([]byte(ORDERLIST))
		c := orders.Cursor()

		// Iterate over all orders, decode and store in the orders list
		for k, v := c.First(); k != nil; k, v = c.Next() {
			order := Order{}
			if err := json.Unmarshal(v, &order); err != nil {
				return err
			}

			orderList = append(orderList, order)
		}

		return
	})

	if err != nil {
		return nil, err
	}

	return &orderList, nil
}
