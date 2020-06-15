package migrater

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestEnv(t *testing.T) {
	if len(os.Getenv("MONGO_HOST")) == 0 {
		t.Fatal("MONGO_HOST env variable is not set properly")
	}
	if len(os.Getenv("MONGO_PORT")) == 0 {
		t.Fatal("MONGO_PORT env variable is not set properly")
	}
}

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

func TestRun(t *testing.T) {
	m := NewMigrater()
	mongoURI := fmt.Sprintf("mongodb://%s:%s", os.Getenv("MONGO_HOST"), os.Getenv("MONGO_PORT"))
	ctx := context.Background()
	clientOpts := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		t.Fatal("Unable to connect to Mongo")
	}
	db := client.Database("migrater")
	m.SetMongoDatabase(db)

	m.Run()
	if m.counter > 0 {
		t.Fatal("Mirations counter should be set to 0")
	}
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

	m.Run()
	if m.counter != 1 {
		t.Fatal("Mongo counter should be set to", "1", "Got", m.counter)
	}
	// clear migrations table
	db.Collection("migrations").DeleteMany(ctx, bson.D{})
}

func TestRunAndIsMigrated(t *testing.T) {
	m := NewMigrater()
	mongoURI := fmt.Sprintf("mongodb://%s:%s", os.Getenv("MONGO_HOST"), os.Getenv("MONGO_PORT"))
	ctx := context.Background()
	clientOpts := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		t.Fatal("Unable to connect to Mongo")
	}
	db := client.Database("migrater")
	m.SetMongoDatabase(db)

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
	m.Run()
	// check migrations count
	count, err := db.Collection("migrations").CountDocuments(ctx, bson.D{})
	if err != nil {
		t.Fatal(err.Error())
	}
	if count != 1 {
		t.Fatal("Documents count in migrations collection should be", "1", "Got", count)
	}
	// add the same migration to fill isMigrated
	m.AddMongoMigration(mig)
	m.Run()
	// check migrations count
	if count != 1 {
		t.Fatal("Documents count in migrations collection should be", "1", "Got", count)
	}
	// clear migrations table
	db.Collection("migrations").DeleteMany(ctx, bson.D{})
}

func TestRollback(t *testing.T) {
	m := NewMigrater()
	mongoURI := fmt.Sprintf("mongodb://%s:%s", os.Getenv("MONGO_HOST"), os.Getenv("MONGO_PORT"))
	ctx := context.Background()
	clientOpts := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		t.Fatal("Unable to connect to Mongo")
	}
	db := client.Database("migrater")
	m.SetMongoDatabase(db)

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
	m.Rollback()
	if m.counter > 0 {
		t.Fatal("Rollback counter should be set to", "1", "Got", m.counter)
	}
	m.AddMongoMigration(mig)
	// run migrations to have some data
	m.Run()
	m.Rollback()
	if m.counter == 0 {
		t.Fatal("Rollback counter should be set to", "1", "Got", m.counter)
	}
	// clear migrations table
	db.Collection("migrations").DeleteMany(ctx, bson.D{})
}
