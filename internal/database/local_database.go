package database

import "gnoty/internal/types"

type localDatabase interface {
	Insert(information types.CommandInformation) error
	List() (*[]types.CommandInformation, error)
	DeleteLocalDb() error
}

func (d databaseService) Insert(information types.CommandInformation) error {
	return d.localDb.Insert(information)
}

func (d databaseService) DeleteLocalDb() error {
	return d.localDb.Delete()
}

func (d databaseService) List() (*[]types.CommandInformation, error) {
	return d.localDb.List()
}
