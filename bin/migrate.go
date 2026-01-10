package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"scopeapi.local/backend/shared/database/postgresql"
)

func main() {
	var (
		host     = flag.String("host", "localhost", "Database host")
		port     = flag.String("port", "5432", "Database port")
		user     = flag.String("user", "postgres", "Database user")
		password = flag.String("password", "password", "Database password")
		dbname   = flag.String("dbname", "scopeapi", "Database name")
		sslmode  = flag.String("sslmode", "disable", "SSL mode")
		action   = flag.String("action", "migrate", "Action: migrate, rollback, status")
	)
	flag.Parse()

	config := postgresql.Config{
		Host:     *host,
		Port:     *port,
		User:     *user,
		Password: *password,
		DBName:   *dbname,
		SSLMode:  *sslmode,
	}

	migrator, err := postgresql.NewMigrator(config)
	if err != nil {
		log.Fatalf("Failed to create migrator: %v", err)
	}
	defer migrator.Close()

	migrationsDir := "backend/shared/database/postgresql/migrations"

	switch *action {
	case "migrate":
		if err := migrator.Migrate(migrationsDir); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		fmt.Println("Migrations completed successfully")
	case "rollback":
		if err := migrator.Rollback(); err != nil {
			log.Fatalf("Rollback failed: %v", err)
		}
		fmt.Println("Rollback completed successfully")
	case "status":
		if err := migrator.Status(migrationsDir); err != nil {
			log.Fatalf("Status check failed: %v", err)
		}
	default:
		fmt.Printf("Unknown action: %s\n", *action)
		fmt.Println("Available actions: migrate, rollback, status")
		os.Exit(1)
	}
}
