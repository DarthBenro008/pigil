package database

type Service interface {
	localDatabase
	configurationDatabase
}

type databaseService struct {
	localDb  BoltDatabase
	configDb BoltDatabase
}

func NewDatabaseService(localDb BoltDatabase, configDb BoltDatabase) Service {
	return &databaseService{localDb: localDb, configDb: configDb}
}
