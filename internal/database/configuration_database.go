package database

import "pigil/internal/types"

type configurationDatabase interface {
	InsertConfig(information types.ConfigurationInformation) error
	GetConfig(key string) (string, error)
	ListConfig() (*[]types.ConfigurationInformation, error)
	DeleteConfigDb() error
	DeleteConfig(key string) error
}

func (d databaseService) InsertConfig(information types.ConfigurationInformation) error {
	return d.configDb.InsertConfig(information)
}

func (d databaseService) GetConfig(key string) (string, error) {
	return d.configDb.GetConfig(key)
}

func (d databaseService) ListConfig() (*[]types.ConfigurationInformation, error) {
	return d.configDb.ListConfig()
}

func (d databaseService) DeleteConfig(key string) error {
	return d.configDb.DeleteConfig(key)
}

func (d databaseService) DeleteConfigDb() error {
	return d.configDb.Delete()
}
