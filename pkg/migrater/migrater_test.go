package migrater

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"

	"bou.ke/monkey"
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
	ctx := context.Background()
	db := connectMongo(t)
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
	collection := db.Collection("migrations")
	// check migrations count
	count, err := collection.CountDocuments(ctx, bson.D{})
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
	collection.DeleteMany(ctx, bson.D{})
}

func TestRunError(t *testing.T) {
	m := NewMigrater()
	ctx := context.Background()
	db := connectMongo(t)
	m.SetMongoDatabase(db)

	mig := MongoMigration{
		Timestamp:   uint64(time.Now().Unix()),
		Description: "Your description",
		Up: func(db *mongo.Database) error {
			return errors.New("Testing purpose error")
		},
		Down: func(db *mongo.Database) error {
			return nil
		},
	}
	m.AddMongoMigration(mig)
	err := m.Run()
	if err == nil {
		t.Error("There should be an error")
	}
	db.Collection("migrations").DeleteMany(ctx, bson.D{})
}

func TestRunSaveMigrationError(t *testing.T) {
	m := NewMigrater()
	ctx := context.Background()
	db := connectMongo(t)
	m.SetMongoDatabase(db)

	// put index on Timestamp to force an error
	collection := db.Collection("migrations")
	collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.M{
				"Description": 1,
			},
			Options: options.Index().SetUnique(true),
		},
	)

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
	err := m.Run()
	if err != nil {
		t.Error(err.Error())
	}
	migWithDescriptionErr := MongoMigration{
		Timestamp:   uint64(20060102150405),
		Description: "Your description",
		Up: func(db *mongo.Database) error {
			return nil
		},
		Down: func(db *mongo.Database) error {
			return nil
		},
	}
	m.AddMongoMigration(migWithDescriptionErr)
	err = m.Run()
	if err == nil {
		t.Error("There should be an error")
	}
	// drop migrations collection
	collection.Drop(ctx)
}

func TestRollback(t *testing.T) {
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
	m.Rollback()
	if m.counter > 0 {
		t.Fatal("Rollback counter should be set to", "1", "Got", m.counter)
	}
	m.AddMongoMigration(mig)
	// run migrations to have some data
	m.Run()
	// reset counter
	m.counter = 0
	m.Rollback()
	if m.counter == 0 {
		t.Fatal("Rollback counter should be set to", "1", "Got", m.counter)
	}
	// clear migrations table
	db.Collection("migrations").DeleteMany(ctx, bson.D{})
}

func TestRollbackError(t *testing.T) {
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
			return errors.New("Testing purpose error")
		},
	}
	m.AddMongoMigration(mig)
	// run migrations to have some data
	m.Run()
	err := m.Rollback()
	if err == nil {
		t.Error("There should be an error")
	}
	// clear migrations table
	db.Collection("migrations").DeleteMany(ctx, bson.D{})
}

func TestRollbackDeleteMigrationError(t *testing.T) {
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
	// run migrations to have some data
	err := m.Run()
	if err != nil {
		t.Error(err.Error())
	}

	// force an error
	var c *mongo.Collection
	var guard *monkey.PatchGuard
	// note that during test there must be flag -gcflags=-l
	guard = monkey.PatchInstanceMethod(reflect.TypeOf(c), "DeleteOne",
		func(c *mongo.Collection, ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
			log.Printf("record: %+v, collection: %s, database: %s", filter, c.Name(), c.Database().Name())
			return nil, errors.New("Test error")
		})
	defer guard.Unpatch()

	err = m.Rollback()
	if err == nil {
		t.Error("There should be an error")
	}
	// clear migrations
	db.Collection("migrations").DeleteMany(ctx, bson.D{})
}

func TestRollbackWithReduce(t *testing.T) {
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
	mig2 := MongoMigration{
		Timestamp:   uint64(time.Now().Unix()),
		Description: "Your description second",
		Up: func(db *mongo.Database) error {
			return nil
		},
		Down: func(db *mongo.Database) error {
			return nil
		},
	}
	m.AddMongoMigration(mig)
	m.AddMongoMigration(mig2)
	// run migrations to have some data
	err := m.Run()
	if err != nil {
		t.Error(err.Error())
	}
	// reset counter
	m.counter = 0
	// pass one argument to Rollback (uint converted to string)
	err = m.Rollback(strconv.FormatUint(mig.Timestamp, 10))
	if err != nil {
		t.Error(err.Error())
	}
	if m.counter != 1 {
		t.Fatal("Rollback counter should be set to", "1", "Got", m.counter)
	}
	// clear migrations
	db.Collection("migrations").DeleteMany(ctx, bson.D{})
}

func TestRollbackWithReduceBadTimestampError(t *testing.T) {
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
	// run migrations to have some data
	err := m.Run()
	if err != nil {
		t.Error(err.Error())
	}
	// reset counter
	m.counter = 0
	// pass one argument to Rollback (uint converted to string)
	err = m.Rollback("bad_timestamp")
	if err == nil {
		t.Error("There should be an error")
	}
	// clear migrations
	db.Collection("migrations").DeleteMany(ctx, bson.D{})
}
