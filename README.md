# Migrater [![codecov](https://codecov.io/gh/malekim/migrater/branch/master/graph/badge.svg)](https://codecov.io/gh/malekim/migrater)
Migrater is a package to easily handle database migrations written in GO

## Install

Simple go get it to your project:
```bash
go get github.com/malekim/migrater
```

To use it as a CLI run:
```bash
go build github.com/malekim/migrater
```

After that a binary with name "migrater" should appear in your project.

## CLI commands

After building binary you can run migrater cli commands.

### Help

```bash
./migrater --help
```

### Generate migration file

```bash
./migrater migration:generate {db_driver}
```

Currently migrater supports only mongodb. To create mongodb migration file use command as follows:

```bash
./migrater migration:generate mongo
```

The above command will generate migration file inside app/migrations.

## Running migrations

To run migrations you have to call similar function:

```go
package yourpackage
import (
  "yourpackage/app/migrations"
  "go.mongodb.org/mongo-driver/mongo"
  "github.com/malekim/migrater/pkg/migrater"
)

//
// some code
//

func RunMigrations(db *mongo.Database) {
  mig := migrater.NewMigrater()
  mig.SetMongoDatabase(db)
  // after creating migration, you have to manually add here a corresponding file
  mig.AddMongoMigration(migrations.Migration1592085513)
  mig.AddMongoMigration(migrations.Migration1592085633)
  err := mig.Run()
  if err != nil {
    // handle err
  }
}
```

## Rollback migrations

To rollback migrations you have to call similar function:

```go
package yourpackage
import (
  "yourpackage/app/migrations"
  "go.mongodb.org/mongo-driver/mongo"
  "github.com/malekim/migrater/pkg/migrater"
)

//
// some code
//

func RollbackMigrations(db *mongo.Database) {
  mig := migrater.NewMigrater()
  mig.SetMongoDatabase(db)
  // after creating migration, you have to manually add here a corresponding file
  mig.AddMongoMigration(migrations.Migration1592085513)
  mig.AddMongoMigration(migrations.Migration1592085633)
  err := mig.Rollback()
  if err != nil {
    // handle err
  }
}
```

To rollback particular migrations, just pass their timestamp as an argument to Rollback function:

```go
func RollbackMigrations(db *mongo.Database) {
  mig := migrater.NewMigrater()
  mig.SetMongoDatabase(db)
  // after creating migration, you have to manually add here a corresponding file
  mig.AddMongoMigration(migrations.Migration1592085513)
  mig.AddMongoMigration(migrations.Migration1592085633)
  // only 1592085513 migration will be reverted
  err := mig.Rollback("1592085513")
  if err != nil {
    // handle err
  }
}
```

The best options is to use it as CLI. Example with cobra:

```go
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run migration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Help()
		return fmt.Errorf("%s requires a subcommand", cmd.Name())
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback migration",
	RunE: func(cmd *cobra.Command, args []string) error {
		db := GetMongoDB()
		mig := migrater.NewMigrater()
		mig.SetMongoDatabase(db)
		mig.AddMongoMigration(migrations.Migration1592085512)
		mig.AddMongoMigration(migrations.Migration1592085633)
		err := mig.Rollback(args...)

		return err
	},
}
```

Then you can call it like:

1. Rolback all migrations:

```bash
go run ./main.go migrate down
```

1. Rollback single migration

```bash
go run ./main.go migrate down 1592085513
```

1. Rollback multiple migrations

```bash
go run ./main.go migrate down 1592085513 1592085633
```

## How it works

When migrater creates a new migration, there are two methods to implement: up and down.

Up method is called during migration. Down method is called during migration rollback.

## Running tests

Ensure that you have working mongo database and pass to test MONGO_HOST and MONGO_PORT:
```bash
MONGO_HOST=localhost MONGO_PORT=27017 go test -v -gcflags=-l -coverprofile=coverage.txt -covermode=atomic ./... &&  go tool cover -html=coverage.txt
```

Note that flag -gcflags=-l is necessary for bou.ke/monkey library.