package main

import (
	"fmt"
	bolt "go.etcd.io/bbolt"
)

type BoltDatabase interface {
	Insert(information CommandInformation) error
	Delete() error
	List() (*[]CommandInformation, error)
}

type boltDatabase struct {
	BoltDb *bolt.DB
}

func (b boltDatabase) Insert(information CommandInformation) error {
	return b.BoltDb.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(Stba(bucket))
		err = b.Put(Stba("asdf"), Stba("asdf"))
		return err
	})
}

func (b boltDatabase) Delete() error {
	return b.BoltDb.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(Stba(bucket))
		if err != nil {
			return err
		}
		return nil
	})
}

func (b boltDatabase) List() (*[]CommandInformation, error) {
	var ci []CommandInformation
	err := b.BoltDb.View(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(Stba(bucket))
		if err != nil {
			return err
		}
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("key=%s, value=%s\n", k, v)
		}
		return nil
	})
	return &ci, err
}

func NewBoldDbService(boltDb *bolt.DB) BoltDatabase {
	return &boltDatabase{BoltDb: boltDb}
}
