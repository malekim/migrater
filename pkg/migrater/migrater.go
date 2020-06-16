package migrater

import (
	"log"
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
		mongo:   &MongoMigrater{},
	}
}

func (m *migrater) AddMongoMigration(mgtn MongoMigration) {
	m.mongo.migrations = append(m.mongo.migrations, mgtn)
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

func (m *migrater) Rollback() error {
	for _, migration := range m.mongo.migrations {
		// check if migration was called before
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
	}
	if m.counter == 0 {
		log.Println("There was nothing to rollback")
	}
	return nil
}
