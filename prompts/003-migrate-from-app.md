# Migrate from app

> I would like to be able to run the migrations from within the app as well, when it starts up

You can integrate the golang-migrate/migrate library directly into your app to run migrations on startup. First, add the library to your project:

```sh
$ go get -u github.com/golang-migrate/migrate/v4
$ go get -u github.com/golang-migrate/migrate/v4/database/postgres
$ go get -u github.com/golang-migrate/migrate/v4/source/file
```

Next, create a function to run the migrations in the internal/database package:

internal/database/migrations.go:

```go
package database

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(conn *postgres.PgxDriver) error {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	migrationsPath := filepath.Join(basepath, "../../db/migrations")

	sourceDriver, err := file.WithPath(migrationsPath)
	if err != nil {
		return fmt.Errorf("unable to create file source driver: %v", err)
	}

	m, err := migrate.NewWithInstance("file", sourceDriver, "postgres", conn)
	if err != nil {
		return fmt.Errorf("unable to create migration instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %v", err)
	}

	log.Println("Migrations applied successfully")

	return nil
}
```

Now, update the New function in the internal/database/database.go to run the migrations after connecting to the database:

```go
package database

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/lib/pq"
)

type DB struct {
	Conn *sql.DB
}

func New(connectionString string) *DB {
	conn, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}

	pgDriver, err := postgres.WithInstance(conn, &postgres.Config{})
	if err != nil {
		panic(fmt.Errorf("unable to create postgres driver instance: %v", err))
	}

	if err := RunMigrations(pgDriver); err != nil {
		panic(fmt.Errorf("failed to run migrations: %v", err))
	}

	return &DB{Conn: conn}
}
```

With these changes, your app will automatically run the migrations whenever it starts up. If there are no new migrations to apply, the app will continue running without any issues.

## Note

The generated code actually had some bugs in `migrations.go`:

- the type `*postgres.PgxDriver` doesn't actually exist, and had to be replaced with `database.Driver` (without the `*`);
- similarly, `file.WithPath` doesn't actually exist. Fortunately, after some back-and-forth, ChatGPT came up with a suitable replacement in a separate conversation:

```go
	migrationUrl := fmt.Sprintf("file://%s", migrationsPath)

	m, err := migrate.NewWithDatabaseInstance(migrationUrl, "postgres", conn)

```
