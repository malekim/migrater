package migrater

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

func TestNewMigrater(t *testing.T) {
	m := NewMigrater()

	mt := reflect.TypeOf(m)

	if fmt.Sprintf("%T", m) != "*migrater.migrater" {
		t.Fatal("Expected", "*migrater.migrater", "Got", mt)
	}

	if m.counter > 0 {
		t.Fatal("Expected", 0, "Got", m.counter)
	}
}

func TestAddMongoMigration(t *testing.T) {
	m := NewMigrater()
	mig := MongoMigration{
		Timestamp:   uint64(time.Now().Unix()),
		Description: "Your description",
		Up: func(db *mongo.Database) error {
			return nil
		},
		Down: func(db *mongo.Database) error {
			return nil
		},
	}
	m.AddMongoMigration(mig)

	if len(m.mongo.migrations) == 0 {
		t.Fatal("Expected", len(m.mongo.migrations), "Got", 0)
	}
}
