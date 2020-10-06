package migrater

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// migrater is basic struct to help
// to make all migrations
//
// counter is set during migration
type migrater struct {
	counter uint
	mongo   *MongoMigrater
	table   string
}

func NewMigrater() *migrater {
	return &migrater{
		counter: 0,
		mongo:   NewMongoMigrater(),
	}
}

func (m *migrater) AddMongoMigration(mgtn MongoMigration) {
	st := strconv.FormatUint(mgtn.Timestamp, 10)
	m.mongo.migrations[st] = mgtn
}

func (m *migrater) SetMongoDatabase(db *mongo.Database) {
	m.mongo.db = db
}

func (m *migrater) Run() error {
	// run mongo migrations
	for _, migration := range m.mongo.migrations {
		// check if migration was called before
		if m.mongo.IsMigrated(migration.Timestamp) {
			continue
		}
		err := migration.Up(m.mongo.db)
		if err != nil {
			return err
		}
		// increment counter
		m.counter++
		// save information about migration to database
		en := &MongoMigrationEntity{
			Timestamp:   migration.Timestamp,
			Description: migration.Description,
			Migrated:    time.Now(),
		}
		err = m.mongo.SaveMigration(en)
		if err != nil {
			return err
		}
		log.Printf("Migration %d (%s) succeded", migration.Timestamp, migration.Description)
	}
	if m.counter == 0 {
		log.Println("There was nothing to migrate")
	}
	return nil
}

func (m *migrater) Rollback(timestamps ...string) error {
	err := m.reduceMigrations(timestamps...)
	if err != nil {
		return err
	}

	for _, migration := range m.mongo.migrations {
		err := m.rollbackOne(migration)
		if err != nil {
			return err
		}
	}

	if m.counter == 0 {
		log.Println("There was nothing to rollback")
	}
	return nil
}

func (m *migrater) rollbackOne(migration MongoMigration) error {
	if m.mongo.IsMigrated(migration.Timestamp) {
		err := migration.Down(m.mongo.db)
		if err != nil {
			return err
		}
		// increment counter
		m.counter++
		err = m.mongo.DeleteMigration(migration.Timestamp)
		if err != nil {
			return err
		}

		log.Printf("Rollback migration %d (%s) succeded", migration.Timestamp, migration.Description)
	}
	return nil
}

func (m *migrater) reduceMigrations(timestamps ...string) error {
	if len(timestamps) == 0 {
		return nil
	}
	new := make(map[string]MongoMigration)

	for _, t := range timestamps {
		if migration, ok := m.mongo.migrations[t]; ok {
			st := strconv.FormatUint(migration.Timestamp, 10)
			new[st] = migration
		} else {
			return fmt.Errorf("Migration with timestamp: `%s` does not exist or has not been added to migrations map.", t)
		}
	}

	m.mongo.migrations = new

	return nil
}
