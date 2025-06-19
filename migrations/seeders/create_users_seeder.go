//go:build ignore

package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/RindangRamadhan/dealls-payslip-system/config"
	"github.com/RindangRamadhan/dealls-payslip-system/pkg/utils"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func main() {
	// Load Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	dbURL := fmt.Sprintf("%s?sslmode=disable", cfg.PG.URL)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("failed to connect to DB:", err)
	}
	defer db.Close()

	if err := SeedData(db); err != nil {
		log.Fatal("seeding failed:", err)
	}
}

func SeedData(db *sql.DB) error {
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM users WHERE is_admin = true`).Scan(&count)
	if err != nil {
		return fmt.Errorf("count check error: %w", err)
	}

	if count > 0 {
		log.Println("Admin user already exists, skipping seeding")
		return nil
	}

	// Create admin
	adminID := uuid.New()
	hashedPassword, err := utils.HashPassword("admin123")
	if err != nil {
		return fmt.Errorf("hash admin password: %w", err)
	}

	_, err = db.Exec(`
        INSERT INTO users (id, username, password, is_admin, salary, created_at, updated_at, created_by, updated_by)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		adminID, "admin", hashedPassword, true, 0, time.Now(), time.Now(), adminID, adminID)
	if err != nil {
		return fmt.Errorf("insert admin: %w", err)
	}

	// Create 100 employees
	for i := 1; i <= 100; i++ {
		employeeID := uuid.New()
		username := fmt.Sprintf("employee%d", i)
		password := fmt.Sprintf("password%d", i)
		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			return fmt.Errorf("hash password for %s: %w", username, err)
		}
		salary := 3000 + float64(i%50)*100

		_, err = db.Exec(`
            INSERT INTO users (id, username, password, is_admin, salary, created_at, updated_at, created_by, updated_by)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			employeeID, username, hashedPassword, false, salary, time.Now(), time.Now(), adminID, adminID)
		if err != nil {
			return fmt.Errorf("insert %s: %w", username, err)
		}

		// Logging Progress
		if i%10 == 0 || i == 100 {
			log.Printf("Seeded %d employees...\n", i)
		}
	}

	log.Println("Seeded 1 admin and 100 employees.")
	return nil
}
