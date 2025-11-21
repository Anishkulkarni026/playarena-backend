package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	var err error

	// Read DSN from environment
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("❌ Error: DB_DSN environment variable not set")
	}

	// ⭐ Log the DSN being used ⭐
	log.Println("👉 Using DSN:", dsn)

	// Open DB connection
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("❌ Failed to open DB connection:", err)
	}

	// Ping DB
	err = DB.Ping()
	if err != nil {
		log.Fatal("❌ Failed to ping DB:", err)
	}

	// ⭐ Check which database is actually selected ⭐
	var dbName string
	err = DB.QueryRow("SELECT DATABASE()").Scan(&dbName)
	if err != nil {
		log.Println("⚠️ Could not detect current database:", err)
	} else {
		log.Println("📌 Connected to database:", dbName)
	}

	log.Println("✅ Successfully connected to MySQL database!")
}
