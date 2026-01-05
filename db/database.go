package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

var Conn *pgx.Conn

func Connect(url string) {
	var err error
	Conn, err = pgx.Connect(context.Background(), url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Successfully connected to the database!")
}

func CreateTable() {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		username TEXT NOT NULL,
		name TEXT,
		email TEXT UNIQUE NOT NULL,
		profilePic TEXT,
		password TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	)`
	_, err := Conn.Exec(context.Background(), query)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	migrations := []string{
		"ALTER TABLE users ADD COLUMN IF NOT EXISTS profilePic TEXT",
		"ALTER TABLE users ADD COLUMN IF NOT EXISTS password TEXT",
		"ALTER TABLE users ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT NOW()",
		"ALTER TABLE users ADD CONSTRAINT users_email_key UNIQUE (email)",
	}

	for _, mig := range migrations {
		Conn.Exec(context.Background(), mig) // Ignore errors for existing constraints
	}
	fmt.Println("Users table ready.")
}
