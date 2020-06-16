package migrater

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connectMongo(t *testing.T) *mongo.Database {
	mongoURI := fmt.Sprintf("mongodb://%s:%s", os.Getenv("MONGO_HOST"), os.Getenv("MONGO_PORT"))
	ctx := context.Background()
	clientOpts := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		t.Fatal("Unable to connect to Mongo")
	}
	return client.Database("migrater")
}

func TestIsMigrated(t *testing.T) {
	m := NewMigrater()
	ctx := context.Background()
	db := connectMongo(t)
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
	isMigrated := m.mongo.IsMigrated(mig.Timestamp)
	if !isMigrated {
		t.Fatal("IsMigrated should return", true, "Got", false)
	}
	// clear migrations table
	db.Collection("migrations").DeleteMany(ctx, bson.D{})
}

func TestSaveMigration(t *testing.T) {
	m := NewMigrater()
	ctx := context.Background()
	db := connectMongo(t)
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
	en := &MongoMigrationEntity{
		Timestamp:   mig.Timestamp,
		Description: mig.Description,
		Migrated:    time.Now(),
	}
	err := m.mongo.SaveMigration(en)
	if err != nil {
		t.Fatal(err.Error())
	}
	collection := db.Collection("migrations")
	count, err := collection.CountDocuments(ctx, bson.M{"timestamp": en.Timestamp})
	if count != 1 {
		t.Fatal("Count migrations should return", 1, "Got", count)
	}
	// clear migrations table
	collection.DeleteMany(ctx, bson.D{})
}

func TestDeleteMigration(t *testing.T) {
	m := NewMigrater()
	ctx := context.Background()
	db := connectMongo(t)
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
	en := &MongoMigrationEntity{
		Timestamp:   mig.Timestamp,
		Description: mig.Description,
		Migrated:    time.Now(),
	}
	err := m.mongo.SaveMigration(en)
	if err != nil {
		t.Fatal(err.Error())
	}
	collection := db.Collection("migrations")
	count, err := collection.CountDocuments(ctx, bson.M{"timestamp": en.Timestamp})
	if count != 1 {
		t.Fatal("Count migrations should return", 1, "Got", count)
	}
	if err != nil {
		t.Fatal(err.Error())
	}
	err = m.mongo.DeleteMigration(en.Timestamp)
	if err != nil {
		t.Fatal(err.Error())
	}
	count, err = collection.CountDocuments(ctx, bson.M{"timestamp": en.Timestamp})
	if count > 0 {
		t.Fatal("Count migrations should return", 0, "Got", count)
	}
	if err != nil {
		t.Fatal(err.Error())
	}
	// clear migrations table
	collection.DeleteMany(ctx, bson.D{})
}

func TestAddMongoMigrationFile(t *testing.T) {
	err := AddMongoMigrationFile()
	if err != nil {
		t.Error("Error during call AddMongoMigrationFile")
	}
	// clean after test
	dir := filepath.Join("app")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Errorf("Dir %s should exist", dir)
	}
	// remove testing dir
	err = os.RemoveAll(dir)
	if err != nil {
		t.Errorf("Unsuccessful clear %s", dir)
	}
}

func TestAddMongoMigrationFileError(t *testing.T) {
	dir := filepath.Join("app")
	// force error and create dir
	// with chmod only to read
	err := os.MkdirAll(dir, 0444)
	if err != nil {
		t.Error("Error during creating the directory")
	}
	err = AddMongoMigrationFile()
	if err == nil {
		t.Error("There should be an error")
	}
	// remove testing dir
	err = os.RemoveAll(dir)
	if err != nil {
		t.Errorf("Unsuccessful clear %s", dir)
	}
}
