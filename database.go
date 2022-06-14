package main

type DatabaseService interface {
	Insert(information CommandInformation) error
	Delete() error
	List() (*[]CommandInformation, error)
}

type databaseService struct {
	boltDatabase BoltDatabase
}

type CommandInformation struct {
	CommandName      string
	CommandArguments string
	ExecutionTime    string
	WasSuccessful    bool
}

func (d databaseService) Insert(information CommandInformation) error {
	return d.boltDatabase.Insert(information)
}

func (d databaseService) Delete() error {
	return d.boltDatabase.Delete()
}

func (d databaseService) List() (*[]CommandInformation, error) {
	return d.boltDatabase.List()
}

func NewDatabaseService(database BoltDatabase) DatabaseService {
	return &databaseService{boltDatabase: database}
}
