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
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "groceries.db"
	}

	email := flag.String("email", "", "user email")
	name := flag.String("name", "", "household member name")
	householdID := flag.Int64("household-id", 1, "household ID")
	flag.StringVar(&dbPath, "db", dbPath, "SQLite database path")
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
	db, err := database.New(ctx, dbPath)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	var userID int64
	err = db.WithTx(ctx, func(tx database.Executor) error {
		result, err := tx.ExecContext(
			ctx,
			`INSERT INTO users (email, password_hash) VALUES (?, ?)`,
			*email,
			string(passwordHash),
		)
		if err != nil {
			return fmt.Errorf("insert user: %w", err)
		}

		userID, err = result.LastInsertId()
		if err != nil {
			return fmt.Errorf("get user ID: %w", err)
		}

		_, err = tx.ExecContext(
			ctx,
			`INSERT INTO members (household_id, user_id, name) VALUES (?, ?, ?)`,
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
