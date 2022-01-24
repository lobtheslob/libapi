package libapi

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	log "github.com/sirupsen/logrus"
)

func handleMigrations(db *sql.DB) error {
	if db == nil {
		log.Fatal("Database is nil")
	}
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Error("an error occurred with migrations1: ", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "mysql", driver)
	if err != nil {
		log.Error("error when migrating: ", err)
	}

	version, _, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.WithField("version", version).Error(err)
		return err
	}

	if err = m.Up(); err != nil {
		log.Error("an error occurred with migrations: ", err)
	}

	//	if err := m.Down(); err != nil {
	//		log.Error("an error occurred with migrations3: ", err)
	//	}

	//// when stupid, reset to specific versions (dirty - figure better way) 
	//if err = m.Force(0); err != nil {
	//	log.Error(err)
	//	return err
	//}

	nversion, _, err := m.Version()
	if err != nil {
		log.Error(err)
		return err
	}

	log.Infof("migrated MySQL DB from version %d to version %d", version, nversion)

	return nil
}