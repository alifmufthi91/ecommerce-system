package pkg

import (
	"database/sql"
	"errors"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	Db   *gorm.DB
	Mock sqlmock.Sqlmock
}

func SetupMockDB() (*Database, error) {
	var (
		db  *sql.DB
		err error
	)
	d := &Database{}
	db, d.Mock, err = sqlmock.New()
	if err != nil {
		return d, errors.New("failed to open mock sql db")
	}
	if db == nil {
		return d, errors.New("mock db is null")
	}
	if d.Mock == nil {
		return d, errors.New("sqlmock is null")
	}
	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	d.Db, err = gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return d, errors.New("Failed to open gorm v2 db")
	}
	if d.Db == nil {
		return d, errors.New("gorm db is null")
	}
	return d, nil
}
