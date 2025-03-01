package postgres

import (
	"bufio"
	"database/sql"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq" // postgres driver
)

func CreateCoinRepository(dbName string, db *sql.DB) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = $1)`
	if err := db.QueryRow(query, dbName).Scan(&exists); err != nil {
		log.Fatal("error checking the existence of the database", err)
	}
	if exists {
		log.Printf("DB %q already exists\n", dbName)
	} else {
		log.Printf("Creating a database %q\n", dbName)
		_, err := db.Exec("CREATE DATABASE " + dbName)
		if err != nil {
			log.Fatal("error when creating the database\n", err)
		}
		log.Printf("database %q was created successfully\n", dbName)
	}
}

func InitCoinTables(db *sql.DB) {
	const op = "database.postgres.InitCoinTables"
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS items (
			name varchar PRIMARY KEY,
			price INT NOT NULL
		);
	`)
	if err != nil {
		log.Fatal(op, err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(255) UNIQUE NOT NULL,
			balance INT NOT NULL DEFAULT 1000)`)
	if err != nil {
		log.Fatal(op, err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS transactions (
			id SERIAL PRIMARY KEY,
			from_user VARCHAR(255) REFERENCES users(username),
			to_user VARCHAR(255) REFERENCES users(username),
			amount INT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)`)
	if err != nil {
		log.Fatal(op, err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS inventory_items (
			id SERIAL PRIMARY KEY,
			user_id INT REFERENCES users(id),
			item_type VARCHAR(255) NOT NULL,
			quantity INT NOT NULL DEFAULT 0);`)
	if err != nil {
		log.Fatal(op, err)
	}

	_, err = db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_inventory_unique
		ON inventory_items (user_id, item_type);
	`)

	if err != nil {
		log.Fatal(op, err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS items(
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) UNIQUE NOT NULL,
			price INT NOT NULL);
	`)
	if err != nil {
		log.Fatal(op, err)
	}

	_, err = db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_name_item_unique
		ON items(name);
	`)
	if err != nil {
		log.Fatal(op, err)
	}
	insertItems(db)
	log.Println("all tables are filled in")
}

func insertItems(db *sql.DB) {
	const op = "database.postgres.insertItems"
	file, err := os.Open("items.csv")
	if err != nil {
		log.Fatal(op, err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), ",")
		_, err = db.Exec(`INSERT INTO items (name, price) 
				VALUES ($1, $2);`,
			fields[0], fields[1])
		if err != nil {
			log.Println(op, err)
			continue
		}
	}
	log.Println("the items table is full")

}
