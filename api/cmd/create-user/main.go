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
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const demoHouseholdID = "019535d9-3df7-79fb-b466-fa907fa17f9e"

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://groceries:groceries@localhost:5432/groceries?sslmode=disable"
	}

	email := flag.String("email", "", "user email")
	name := flag.String("name", "", "user display name")
	householdID := flag.String("household-id", demoHouseholdID, "household UUID")
	flag.StringVar(&databaseURL, "database-url", databaseURL, "PostgreSQL connection URL")
	flag.Parse()

	*email = strings.TrimSpace(*email)
	*name = strings.TrimSpace(*name)
	*householdID = strings.TrimSpace(*householdID)
	if *email == "" || *name == "" {
		log.Fatal("email and name are required")
	}
	if *householdID == "" {
		log.Fatal("household-id is required")
	}
	householdUUID, err := uuid.Parse(*householdID)
	if err != nil {
		log.Fatalf("invalid household-id: %v", err)
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

	var userID uuid.UUID
	err = db.GetContext(ctx, &userID,
		`INSERT INTO users (email, password_hash, name, household_id) VALUES ($1, $2, $3, $4) RETURNING id`,
		*email, string(passwordHash), *name, householdUUID,
	)
	if err != nil {
		log.Fatalf("create user: %v", err)
	}

	fmt.Printf("Created user %s for household %s.\n", userID, householdUUID)
}

func readLine(reader *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	value, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("read input: %v", err)
	}

	return strings.TrimRight(value, "\r\n")
}
