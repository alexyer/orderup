package main

import (
	"encoding/json"
	"errors"

	"github.com/boltdb/bolt"
)

// DB calls.
// Available for Slack and REST APIs.

// Get orders for the queue.
// General function to get list of orders from the queue.
// bucket can be ORDERLIST or HISTORY.
func (o *Orderup) getOrderList(queue []byte, bucket []byte) (*[]Order, error) {
	var orderList []Order

	err := o.db.View(func(tx *bolt.Tx) (err error) {
		// Get bucket with restaurants.
		b := tx.Bucket([]byte(QUEUES))

		restaurant := b.Bucket(queue)
		if restaurant == nil {
			return NonExistentQueue(string(queue))
		}

		orders := restaurant.Bucket(bucket)
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

// Get pending orders list for the queue.
func (o *Orderup) getPendingOrderList(queue []byte) (*[]Order, error) {
	return o.getOrderList(queue, []byte(ORDERLIST))
}

// Get history list for the queue.
func (o *Orderup) getHistoryList(queue []byte) (*[]Order, error) {
	return o.getOrderList(queue, []byte(HISTORY))

}

// Create new queue <name>.
func (o *Orderup) createQueue(name []byte) error {
	return o.db.Update(func(tx *bolt.Tx) (err error) {
		// Get bucket with restaurants.
		b := tx.Bucket([]byte(QUEUES))

		// Create new bucket for the new queue.
		r, err := b.CreateBucket(name)
		if err != nil {
			return errors.New("Queue already exists.")
		}

		// Create 2 subbucktes.
		// One for pending orders, another for finished orders.
		_, err = r.CreateBucket([]byte(ORDERLIST))
		_, err = r.CreateBucket([]byte(HISTORY))

		return err
	})
}

// Delete queue <name>.
func (o *Orderup) deleateQueue(name []byte) error {
	return o.db.Update(func(tx *bolt.Tx) (err error) {
		// Get bucket with restaurants.
		b := tx.Bucket([]byte(RESTAURANTS))

		// Create new bucket for the new restaurant.
		if err = b.DeleteBucket(name); err != nil {
			return errors.New("Restaurant does not exist.")
		}

		return err
	})
}
