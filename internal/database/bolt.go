package database

import (
	"encoding/json"
	"github.com/DarthBenro008/pigil/internal/types"
	"github.com/DarthBenro008/pigil/internal/utils"
	bolt "go.etcd.io/bbolt"
	"strconv"
)

type BoltDatabase interface {
	Insert(information types.CommandInformation) error
	Delete() error
	List() (*[]types.CommandInformation, error)
	InsertConfig(information types.ConfigurationInformation) error
	GetConfig(key string) (string, error)
	ListConfig() (*[]types.ConfigurationInformation, error)
	DeleteConfig(key string) error
}

type boltDatabase struct {
	BoltDb     *bolt.DB
	bucketName string
}

func (b boltDatabase) insert(k []byte, v interface{}) error {
	buf, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return b.BoltDb.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(utils.Stba(b.bucketName))
		err = b.Put(k, buf)
		return err
	})
}

func (b boltDatabase) InsertConfig(information types.ConfigurationInformation) error {
	return b.insert(utils.Stba(information.Key), information)
}

func (b boltDatabase) GetConfig(key string) (string, error) {
	var tci types.ConfigurationInformation
	err := b.BoltDb.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(utils.Stba(b.bucketName))
		res := b.Get(utils.Stba(key))
		if len(res) == 0 {
			tci.Value = ""
			return nil
		}
		err = json.Unmarshal(res, &tci)
		return err
	})
	return tci.Value, err
}

func (b boltDatabase) ListConfig() (*[]types.ConfigurationInformation, error) {
	var ci []types.ConfigurationInformation
	err := b.BoltDb.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(utils.Stba(b.bucketName))
		if err != nil {
			return err
		}
		c := b.Cursor()
		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			var tci types.ConfigurationInformation
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

func (b boltDatabase) DeleteConfig(key string) error {
	return b.BoltDb.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(utils.Stba(b.bucketName))
		err = b.Delete(utils.Stba(key))
		return err
	})
}

func (b boltDatabase) Insert(information types.CommandInformation) error {
	return b.insert(utils.Stba(strconv.Itoa(int(information.ExecutionTime))),
		information)
}

func (b boltDatabase) Delete() error {
	return b.BoltDb.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(utils.Stba(b.bucketName))
		if err != nil {
			return err
		}
		return nil
	})
}

func (b boltDatabase) List() (*[]types.CommandInformation, error) {
	var ci []types.CommandInformation
	err := b.BoltDb.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(utils.Stba(b.bucketName))
		if err != nil {
			return err
		}
		c := b.Cursor()
		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			var tci types.CommandInformation
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

func NewBoltDbService(boltDb *bolt.DB, bucketName string) BoltDatabase {
	return &boltDatabase{BoltDb: boltDb, bucketName: bucketName}
}
