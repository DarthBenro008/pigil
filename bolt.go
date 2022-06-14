package main

import (
	"encoding/json"
	bolt "go.etcd.io/bbolt"
	"strconv"
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
	buf, err := json.Marshal(information)
	if err != nil {
		return err
	}
	return b.BoltDb.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(Stba(bucket))
		err = b.Put(Stba(strconv.Itoa(int(information.ExecutionTime))), buf)
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
	err := b.BoltDb.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(Stba(bucket))
		if err != nil {
			return err
		}
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var tci CommandInformation
			err := json.Unmarshal(v, &tci)
			if err != nil {
				return err
			}
			ci = append(ci, tci)
		}
		return nil
	})
	return &ci, err
}

func NewBoldDbService(boltDb *bolt.DB) BoltDatabase {
	return &boltDatabase{BoltDb: boltDb}
}
