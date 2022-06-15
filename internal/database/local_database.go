package database

import "gnoty/internal/types"

type localDatabase interface {
	Insert(information types.CommandInformation) error
	Delete() error
	List() (*[]types.CommandInformation, error)
}

func (d databaseService) Insert(information types.CommandInformation) error {
	return d.localDb.Insert(information)
}

func (d databaseService) Delete() error {
	return d.localDb.Delete()
}

func (d databaseService) List() (*[]types.CommandInformation, error) {
	return d.localDb.List()
}
