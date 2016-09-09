package main

import (
	"errors"
	"fmt"

	"github.com/boltdb/bolt"
)

// OriginService provides operations on data from originating alert
type OriginService interface {
	ProcessAlert(string) (string, error)
}

type originService struct{}

func (originService) ProcessAlert(s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}

	var alerts = []byte("alerts")

	// open boltdb; create if it does not exist
	db, err := bolt.Open("alerts.db", 0600, nil)
	if err != nil {
		return s, ErrDatabaseNotOpen
	}
	defer db.Close()

	// asign kv pair
	key := []byte(s)
	value := []byte("")

	// create bucket
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(alerts)
		if err != nil {
			return ErrBucketCreationFailed
		}

		return nil
	})

	var key_exists = false

	// check if key exists
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(alerts)
		if bucket == nil {
			return ErrBucketNotFound
		}

		val := bucket.Get(key)
		if val != nil {
			key_exists = true
		}

		return nil
	})

	if err != nil {
		return s, err
	}

	if key_exists == true {
		fmt.Printf("key %s exists; doing noting\n", string(key))
		return s, nil
	}

	// store new key
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(alerts)
		if err != nil {
			return ErrBucketCreationFailed
		}

		err = bucket.Put(key, value)
		if err != nil {
			return ErrDBPut
		}

		fmt.Printf("key %s added\n", string(key))

		return nil
	})

	return s, nil
}

// ErrEmpty is returned when an input string is empty.
var ErrEmpty = errors.New("empty string")

var (

	// ErrBucketCreationFailed
	ErrBucketCreationFailed = errors.New("bucket creation failed")

	// ErrDatabaseNotOpen is returned when a DB instance is accessed before it
	// is opened or after it is closed.
	ErrDatabaseNotOpen = errors.New("database not open")

	// ErrBucketNotFound is returned when trying to access a bucket that has
	// not been created yet.
	ErrBucketNotFound = errors.New("bucket not found")

	// ErrDBPut is returned when there is an issue with adding kv
	ErrDBPut = errors.New("error adding kv")
)
