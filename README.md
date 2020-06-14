# Migrater
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
  mig.Run()
}
```

## How it works

When migrater creates a new migration, there are two methods to implement: up and down.

Up method is called during migration. Down method is called during migrations rollback.

## Running tests

```bash
go test ./...
```