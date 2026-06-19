package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer conn.Close(ctx)

	var id string
	var email string
	var passHash string
	err = conn.QueryRow(ctx, "SELECT id, email, password_hash FROM staffs WHERE email = 'admin@salon-core.com'").Scan(&id, &email, &passHash)
	if err != nil {
		fmt.Printf("User not found: %v\n", err)
		return
	}

	fmt.Printf("User found: %s (%s)\n", id, email)

	err = bcrypt.CompareHashAndPassword([]byte(passHash), []byte("admin123"))
	if err != nil {
		fmt.Printf("Password DOES NOT match admin123!\n")
	} else {
		fmt.Printf("Password matches admin123!\n")
	}
}
