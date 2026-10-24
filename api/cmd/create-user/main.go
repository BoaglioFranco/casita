package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/BoaglioFranco/groceriesList/api/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://groceries:groceries@localhost:5432/groceries?sslmode=disable"
	}

	email := flag.String("email", "", "user email")
	name := flag.String("name", "", "household member name")
	householdID := flag.Int64("household-id", 1, "household ID")
	flag.StringVar(&databaseURL, "database-url", databaseURL, "PostgreSQL connection URL")
	flag.Parse()

	*email = strings.TrimSpace(*email)
	*name = strings.TrimSpace(*name)
	if *email == "" || *name == "" {
		log.Fatal("email and name are required")
	}
	if *householdID <= 0 {
		log.Fatal("household-id must be greater than zero")
	}

	reader := bufio.NewReader(os.Stdin)
	password := readLine(reader, "Password: ")
	confirmation := readLine(reader, "Confirm password: ")
	if password == "" {
		log.Fatal("password is required")
	}
	if password != confirmation {
		log.Fatal("passwords do not match")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("hash password: %v", err)
	}

	ctx := context.Background()
	db, err := database.New(ctx, databaseURL)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	var userID int64
	err = db.WithTx(ctx, func(tx database.Executor) error {
		err := tx.GetContext(
			ctx,
			&userID,
			`INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id`,
			*email,
			string(passwordHash),
		)
		if err != nil {
			return fmt.Errorf("insert user: %w", err)
		}

		_, err = tx.ExecContext(
			ctx,
			`INSERT INTO members (household_id, user_id, name) VALUES ($1, $2, $3)`,
			*householdID,
			userID,
			*name,
		)
		if err != nil {
			return fmt.Errorf("insert household member: %w", err)
		}

		return nil
	})
	if err != nil {
		log.Fatalf("create user: %v", err)
	}

	fmt.Printf("Created user %d for household %d.\n", userID, *householdID)
}

func readLine(reader *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	value, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("read input: %v", err)
	}

	return strings.TrimRight(value, "\r\n")
}
