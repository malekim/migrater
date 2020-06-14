package migrater

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var MongoStub string = `
package migrations
import (
	"github.com/malekim/migrater/pkg/migrater"
	"go.mongodb.org/mongo-driver/mongo"
)

var Migration{{ .Timestamp }} migrater.MongoMigration = migrater.MongoMigration{
	Timestamp:     {{ .Timestamp }},
	Description: "Your description",
	Up: func(db *mongo.Database) error {
		return nil
	},
	Down: func(db *mongo.Database) error {
		return nil
	},
}
`

type MongoMigrater struct {
	counter    uint
	migrations []MongoMigration
	db         *mongo.Database
}

type MongoMigrationFunc func(db *mongo.Database) error

type MongoMigration struct {
	Timestamp   uint64
	Description string
	Up          MongoMigrationFunc
	Down        MongoMigrationFunc
}

type MongoMigrationEntity struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Timestamp   uint64             `json:"timestamp" bson:"timestamp"`
	Description string             `json:"description" bson:"description"`
	Migrated    time.Time          `json:"migrated" bson:"migrated"`
}

func (mgo *MongoMigrater) IsMigrated(timestamp uint64) bool {
	en := &MongoMigrationEntity{}
	collection := mgo.db.Collection("migrations")
	err := collection.FindOne(context.TODO(), bson.M{"timestamp": timestamp}).Decode(&en)
	if err != nil {
		return false
	}
	return true
}

func (mgo *MongoMigrater) SaveMigration(en *MongoMigrationEntity) error {
	collection := mgo.db.Collection("migrations")
	_, err := collection.InsertOne(context.TODO(), en)
	if err != nil {
		return err
	}
	return nil
}

func (mgo *MongoMigrater) DeleteMigration(timestamp uint64) error {
	collection := mgo.db.Collection("migrations")
	_, err := collection.DeleteOne(context.TODO(), bson.M{"timestamp": timestamp})
	if err != nil {
		return err
	}
	return nil
}
