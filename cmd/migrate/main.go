package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/neokofg/go-test/database"
	"log"
	"os"
	"time"
)

func main() {
	cmd := flag.String("command", "", "migrate command (up/down/goto/create)")
	name := flag.String("name", "", "name for create command")
	version := flag.Uint("version", 0, "migration version for goto command")
	flag.Parse()
	connStr := "host=localhost port=5432 user=postgres password=password dbname=sweetify sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	migrator, err := database.NewMigrator(db, "database/migrations")
	if err != nil {
		log.Fatal(err)
	}

	switch *cmd {
	case "up":
		err = migrator.Up()
	case "down":
		err = migrator.Down()
	case "goto":
		err = migrator.Goto(*version)
	case "create":
		err = CreateMigration(*name)
	default:
		log.Fatal("Unknown command. Use up, down, or goto")
	}

	if err != nil {
		log.Fatal(err)
	}
}

func CreateMigration(name string) error {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	baseName := fmt.Sprintf("%s_%s", timestamp, name)

	// Создаём up миграцию
	upFile := fmt.Sprintf("database/migrations/%s.up.sql", baseName)
	if err := os.WriteFile(upFile, []byte("-- Миграция вверх"), 0644); err != nil {
		return err
	}

	// Создаём down миграцию
	downFile := fmt.Sprintf("database/migrations/%s.down.sql", baseName)
	if err := os.WriteFile(downFile, []byte("-- Миграция вниз"), 0644); err != nil {
		return err
	}

	return nil
}
