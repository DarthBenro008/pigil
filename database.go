package main

type DatabaseService interface {
	Open() error
	Insert() error
	Delete() error
	List() error
}

type CommandInformation struct {
	CommandName      string
	CommandArguments string
	ExecutionTime    string
	WasSuccessful    bool
}
