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

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("❌ Error: DB_DSN environment variable not set")
	}

	// ⭐ ADD THIS LINE ⭐
	log.Println("👉 Using DSN:", dsn)

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("❌ Failed to open DB connection:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("❌ Failed to ping DB:", err)
	}

	log.Println("✅ Successfully connected to MySQL database!")
}
