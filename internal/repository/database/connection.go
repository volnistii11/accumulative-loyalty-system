package database

import (
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConnection(driver string, dsn string) (*gorm.DB, error) {
	var (
		db  *gorm.DB
		err error
	)

	if driver != "pgx" {
		errors.New("service only supports postgreSQL now")
	}

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		errors.Wrap(err, "failed to create a db connection")
		return nil, err
	}

	dbInstance, _ := db.DB()
	if err = dbInstance.Ping(); err != nil {
		errors.Wrap(err, "failed to ping the db")
		return nil, err
	}

	return db, nil
}
