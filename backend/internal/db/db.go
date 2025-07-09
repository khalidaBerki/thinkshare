package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
)

// DB est la connexion globale SQL
var DB *sql.DB

// GormDB est la connexion globale GORM
var GormDB *gorm.DB

// InitDB initialise la connexion PostgreSQL et GORM
func InitDB() {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("PGHOST"),
		os.Getenv("PGPORT"),
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
		os.Getenv("PGDATABASE"),
		os.Getenv("PGSSLMODE"),
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Erreur ouverture DB : %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Impossible de ping la DB : %v", err)
	}

	// GORM avec la DSN, pas avec Conn: DB
	GormDB, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatalf("Erreur ouverture GORM : %v", err)
	}

	log.Println("✅ Connexion PostgreSQL & GORM OK")
}
