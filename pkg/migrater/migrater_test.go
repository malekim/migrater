package migrater

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func TestSetMongoDatabase(t *testing.T) {
	m := NewMigrater()
	mongoURI := fmt.Sprintf("mongodb://%s:%s", os.Getenv("MONGO_HOST"), os.Getenv("MONGO_PORT"))
	ctx := context.Background()
	clientOpts := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		t.Fatal("Unable to connect to Mongo")
	}
	if m.mongo.db != nil {
		t.Fatal("Mongo should not be set")
	}
	db := client.Database("migrater")
	m.SetMongoDatabase(db)

	if m.mongo.db == nil {
		t.Fatal("Mongo is not set properly")
	}
}
