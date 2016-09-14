package main

import (
	"errors"
	"os"

	"github.com/boltdb/bolt"
	"github.com/go-kit/kit/log"
)

// OriginService provides operations on data from originating alert
type OriginService interface {
	ProcessAlert(string) (string, error)
	List(string) (map[string]interface{}, error)
}

type originService struct{}

func (originService) ProcessAlert(s string) (string, error) {
	logger := log.NewLogfmtLogger(os.Stderr)
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
		logger.Log("key", s, "exists", true, "action", nil)
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

		logger.Log("key", s, "exist", false, "action", "added")
		return nil
	})

	return s, nil
}

func (originService) List(s string) (map[string]interface{}, error) {
	if s == "" {
		return nil, ErrEmpty
	}

	logger := log.NewLogfmtLogger(os.Stderr)
	var alerts = []byte("alerts")

	logger.Log("method", "list")

	// open boltdb; create if it does not exist
	db, err := bolt.Open("alerts.db", 0600, nil)
	if err != nil {
		return nil, ErrDatabaseNotOpen
	}
	defer db.Close()

	list := make(map[string]interface{})

	// iterate over keys
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(alerts)
		if bucket == nil {
			return ErrBucketNotFound
		}

		err = bucket.ForEach(func(k, v []byte) error {
			list[string(k)] = string(v)
			return nil
		})

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return list, nil
}

var (
	ErrEmpty                = errors.New("empty string")
	ErrBucketCreationFailed = errors.New("bucket creation failed")
	ErrDatabaseNotOpen      = errors.New("database not open")
	ErrBucketNotFound       = errors.New("bucket not found")
	ErrDBPut                = errors.New("error adding kv")
)
