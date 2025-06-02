package Database

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var (
	Todo *sqlx.DB
)

func migrateUp(db *sqlx.DB) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://Database/migrations",
		"postgres", driver)

	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
func Database() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbUrl := os.Getenv("DB_URL")
	connStr := fmt.Sprintf("%s?sslmode=require", dbUrl)
	DB, err := sqlx.Open("postgres", connStr)
	if err != nil {
		log.Fatalln(err)
	}

	err = DB.Ping()
	if err != nil {
		return err
	}

	Todo = DB
	return migrateUp(DB)
}
