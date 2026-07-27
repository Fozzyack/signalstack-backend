package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/Fozzyack/signalstack-backend/internal/database"
	"github.com/Fozzyack/signalstack-backend/migrations"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type seedUser struct {
	Email string
	Name  string
}

type seedRequest struct {
	Title       string
	Description string
	ClientName  string
	ClientEmail string
	Status      string
}

func main() {
	_ = godotenv.Load()

	db, err := database.Open()
	if err != nil {
		fatal("open database", err)
	}
	defer db.Close()

	if err := database.Migrate(db, migrations.FS); err != nil {
		fatal("run migrations", err)
	}

	if err := seed(db); err != nil {
		fatal("seed database", err)
	}

	fmt.Println("database seeded successfully")
	fmt.Println("seed user password:", seedPassword())
}

func seed(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(seedPassword()), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash seed password: %w", err)
	}

	users := []seedUser{
		{Email: "maya.chen@signalstack.test", Name: "Maya Chen"},
		{Email: "james.doyle@signalstack.test", Name: "James Doyle"},
		{Email: "rina.kapoor@signalstack.test", Name: "Rina Kapoor"},
		{Email: "alex.lee@signalstack.test", Name: "Alex Lee"},
		{Email: "priya.shah@signalstack.test", Name: "Priya Shah"},
	}
	userIDs := make(map[string]string, len(users))
	for _, user := range users {
		id, err := upsertUser(tx, user, string(passwordHash))
		if err != nil {
			return err
		}
		userIDs[user.Email] = id
	}

	requests := []seedRequest{
		{
			Title:       "VPN access dropping every 20 minutes",
			Description: "Remote warehouse team is losing access to internal tools during shifts.",
			ClientName:  "Northstar Logistics",
			ClientEmail: "maya.chen@northstar.io",
			Status:      "new",
		},
		{
			Title:       "New starter account and laptop setup",
			Description: "Three designers joining Monday need accounts, devices, and shared drive access.",
			ClientName:  "Arc & Field Studio",
			ClientEmail: "oliver@arcandfield.co",
			Status:      "new",
		},
		{
			Title:       "Review permissions for finance shared drive",
			Description: "Quarterly access review requested before the external audit begins.",
			ClientName:  "Verity & Co",
			ClientEmail: "sarah@verityandco.com",
			Status:      "in_progress",
		},
		{
			Title:       "Move production database to managed cloud",
			Description: "Looking for a migration plan and someone to own the first phase.",
			ClientName:  "Kiteworks",
			ClientEmail: "devops@kiteworks.dev",
			Status:      "waiting",
		},
		{
			Title:       "Harden administrator accounts with MFA",
			Description: "The operations team needs help rolling out stronger sign-in controls for privileged accounts.",
			ClientName:  "Granite Health",
			ClientEmail: "it@granitehealth.example",
			Status:      "in_progress",
		},
		{
			Title:       "Investigate suspicious mailbox forwarding rule",
			Description: "A mailbox has an unexpected forwarding rule and needs an incident review.",
			ClientName:  "Brightwell Partners",
			ClientEmail: "security@brightwell.example",
			Status:      "new",
		},
	}
	requestIDs := make(map[string]string, len(requests))
	for _, request := range requests {
		id, err := upsertRequest(tx, request)
		if err != nil {
			return err
		}
		requestIDs[request.Title] = id
	}

	assignments := []struct {
		requestTitle string
		userEmail    string
		role         string
	}{
		{requestTitle: "VPN access dropping every 20 minutes", userEmail: "maya.chen@signalstack.test", role: "lead"},
		{requestTitle: "Review permissions for finance shared drive", userEmail: "james.doyle@signalstack.test", role: "lead"},
		{requestTitle: "Review permissions for finance shared drive", userEmail: "rina.kapoor@signalstack.test", role: "contributor"},
		{requestTitle: "Move production database to managed cloud", userEmail: "alex.lee@signalstack.test", role: "lead"},
		{requestTitle: "Harden administrator accounts with MFA", userEmail: "priya.shah@signalstack.test", role: "lead"},
		{requestTitle: "Investigate suspicious mailbox forwarding rule", userEmail: "priya.shah@signalstack.test", role: "contributor"},
	}
	for _, assignment := range assignments {
		if err := upsertAssignment(
			tx,
			requestIDs[assignment.requestTitle],
			userIDs[assignment.userEmail],
			assignment.role,
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func upsertUser(tx *sql.Tx, user seedUser, passwordHash string) (string, error) {
	var id string
	err := tx.QueryRow(`
		INSERT INTO users (email, name, password_hash)
		VALUES ($1, $2, $3)
		ON CONFLICT (email) DO UPDATE SET name = EXCLUDED.name
		RETURNING id
	`, user.Email, user.Name, passwordHash).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("upsert user %s: %w", user.Email, err)
	}
	return id, nil
}

func upsertRequest(tx *sql.Tx, request seedRequest) (string, error) {
	var id string
	err := tx.QueryRow(`
		UPDATE requests
		SET description = $1,
			client_name = $2,
			status = $3,
			updated_at = NOW()
		WHERE title = $4 AND client_email = $5
		RETURNING id
	`, request.Description, request.ClientName, request.Status, request.Title, request.ClientEmail).Scan(&id)
	if err == nil {
		return id, nil
	}
	if err != sql.ErrNoRows {
		return "", fmt.Errorf("update seed request %s: %w", request.Title, err)
	}

	err = tx.QueryRow(`
		INSERT INTO requests (title, description, client_name, client_email, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, request.Title, request.Description, request.ClientName, request.ClientEmail, request.Status).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("insert seed request %s: %w", request.Title, err)
	}
	return id, nil
}

func upsertAssignment(tx *sql.Tx, requestID, userID, role string) error {
	result, err := tx.Exec(`
		UPDATE request_assignments
		SET role = $1
		WHERE request_id = $2 AND user_id = $3 AND unassigned_at IS NULL
	`, role, requestID, userID)
	if err != nil {
		return fmt.Errorf("upsert assignment: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check assignment update: %w", err)
	}
	if rows > 0 {
		return nil
	}

	_, err = tx.Exec(`
		INSERT INTO request_assignments (request_id, user_id, role, assigned_at)
		VALUES ($1, $2, $3, $4)
	`, requestID, userID, role, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("insert assignment: %w", err)
	}
	return nil
}

func seedPassword() string {
	if password := os.Getenv("SEED_PASSWORD"); password != "" {
		return password
	}
	return "signalstack-dev"
}

func fatal(action string, err error) {
	fmt.Fprintf(os.Stderr, "failed to %s: %v\n", action, err)
	os.Exit(1)
}
