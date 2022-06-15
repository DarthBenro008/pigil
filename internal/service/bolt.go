package service

import (
	"encoding/json"
	"gnoty/internal/database"
	"gnoty/internal/utils"
	bolt "go.etcd.io/bbolt"
	"strconv"
)

type BoltDatabase interface {
	Insert(information database.CommandInformation) error
	Delete() error
	List() (*[]database.CommandInformation, error)
}

type boltDatabase struct {
	BoltDb *bolt.DB
}

func (b boltDatabase) Insert(information database.CommandInformation) error {
	buf, err := json.Marshal(information)
	if err != nil {
		return err
	}
	return b.BoltDb.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(utils.Stba(utils.Bucket))
		err = b.Put(utils.Stba(strconv.Itoa(int(information.ExecutionTime))), buf)
		return err
	})
}

func (b boltDatabase) Delete() error {
	return b.BoltDb.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(utils.Stba(utils.Bucket))
		if err != nil {
			return err
		}
		return nil
	})
}

func (b boltDatabase) List() (*[]database.CommandInformation, error) {
	var ci []database.CommandInformation
	err := b.BoltDb.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(utils.Stba(utils.Bucket))
		if err != nil {
			return err
		}
		c := b.Cursor()
		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			var tci database.CommandInformation
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
