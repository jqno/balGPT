# Initialize go project

> We're going to build balGPT, an app that tries to predict soccer matches. We're going to use Golang and Postgres for the backend, and we'll figure out the frontend later. We'll deploy to fly.io but we'll also run locally with Postgres in a docker container. Let's start with the scaffolding.

Great! Let's start by setting up the project structure, initializing the Go modules, and creating the necessary files and folders for your balGPT app.

1. Create the project folder and navigate to it:

```shell
$ mkdir balGPT
$ cd balGPT
```

1. Initialize Go modules:

```shell
$ go mod init github.com/yourusername/balGPT
```

1. Create the folder structure:

```shell
$ mkdir -p cmd/nacgpt-api
$ mkdir -p internal/app
$ mkdir -p internal/config
$ mkdir -p internal/database
$ mkdir -p internal/predictor
```

1. Create the main Go file for your API server:

```shell
$ touch cmd/nacgpt-api/main.go
```

1. Create the configuration, application, and database files:

```shell
$ touch internal/config/config.go
$ touch internal/app/app.go
$ touch internal/database/database.go
$ touch internal/predictor/predictor.go
```

Now, let's write some code.

cmd/nacgpt-api/main.go - This file will be the entry point of your API server:

```go
package main

import (
	"github.com/yourusername/balGPT/internal/app"
	"github.com/yourusername/balGPT/internal/config"
)

func main() {
	cfg := config.LoadConfig()
	app := app.NewApp(cfg)
	app.Run()
}
```

internal/config/config.go - This file will handle loading the configuration:

```go
package config

import (
	"os"
)

type Config struct {
	DBConnectionString string
}

func LoadConfig() *Config {
	return &Config{
		DBConnectionString: os.Getenv("DB_CONNECTION_STRING"),
	}
}
```

internal/app/app.go - This file will handle creating and running the application:

```go
package app

import (
	"github.com/yourusername/balGPT/internal/config"
	"github.com/yourusername/balGPT/internal/database"
)

type App struct {
	Config *config.Config
	DB     *database.DB
}

func NewApp(cfg *config.Config) *App {
	db := database.New(cfg.DBConnectionString)
	return &App{
		Config: cfg,
		DB:     db,
	}
}

func (a *App) Run() {
	// TODO: Initialize and run your web server here
}
```

internal/database/database.go - This file will handle connecting to your Postgres database:

```go
package database

import (
	"database/sql"
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
	return &DB{Conn: conn}
}
```

internal/predictor/predictor.go - This file will hold the logic for predicting soccer matches:

```go
package predictor

// TODO: Implement your match prediction logic here
```

Now that you have the basic scaffolding in place, you can proceed with implementing the match prediction logic, setting up the API endpoints, and creating the frontend.

For running the Postgres database locally using Docker, you can use the following command:

```css
$ docker run --name nacgpt-postgres -e POSTGRES_PASSWORD=yourpassword -p 5432:5432 -d postgres
```

To deploy your app to fly.io, you'll need to follow their official guide and create a Dockerfile and .fly/fly.toml file for your project.
