package migrater

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/malekim/migrater/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var mongoStub string = `
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
	migrations map[string]MongoMigration
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

func NewMongoMigrater() *MongoMigrater {
	return &MongoMigrater{
		counter:    0,
		migrations: make(map[string]MongoMigration),
	}
}

func (mgo *MongoMigrater) IsMigrated(timestamp uint64) bool {
	en := &MongoMigrationEntity{}
	collection := mgo.db.Collection("migrations")
	err := collection.FindOne(context.TODO(), bson.M{"timestamp": timestamp}).Decode(&en)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

func (mgo *MongoMigrater) SaveMigration(en *MongoMigrationEntity) error {
	collection := mgo.db.Collection("migrations")
	_, err := collection.InsertOne(context.TODO(), en)
	return err
}

func (mgo *MongoMigrater) DeleteMigration(timestamp uint64) error {
	collection := mgo.db.Collection("migrations")
	_, err := collection.DeleteOne(context.TODO(), bson.M{"timestamp": timestamp})
	return err
}

func AddMongoMigrationFile() error {
	timestamp := time.Now().Unix()
	name := fmt.Sprintf("%d.go", timestamp)
	t := template.Must(template.New("").Parse(mongoStub))
	path := filepath.Join("app", "migrations", name)
	if err := utils.EnsureDir(path); err != nil {
		log.Printf("Error creating dir: %s", err.Error())
		return err
	}
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Error opening file: %s", err.Error())
		return err
	}
	defer f.Close()

	vars := struct {
		Timestamp int64
	}{
		timestamp,
	}

	err = t.Execute(f, vars)
	if err == nil {
		log.Printf("Created %s\n", name)
	}
	return err
}
