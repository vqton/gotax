package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	seedUsername = "admin"
	seedPassword = "Admin@123456!"
	seedFullName = "System Admin"
	seedEmail    = "admin@gotax.vn"
	seedRole     = "admin"
)

func main() {
	dsn := flag.String("dsn", os.Getenv("DATABASE_URL"), "PostgreSQL DSN (env DATABASE_URL)")
	username := flag.String("username", seedUsername, "Login name")
	password := flag.String("password", seedPassword, "Password (>= 12 chars)")
	fullName := flag.String("full-name", seedFullName, "Full name")
	email := flag.String("email", seedEmail, "Email")
	role := flag.String("role", seedRole, "Role: admin | chief_accountant | accountant | viewer")
	force := flag.Bool("force", false, "Overwrite existing user password")
	flag.Parse()

	if *dsn == "" {
		log.Fatal("DATABASE_URL required (env or -dsn flag)")
	}
	if err := validateRole(*role); err != nil {
		log.Fatalf("bad role: %v", err)
	}
	if len(*password) < 12 {
		log.Fatal("password must be at least 12 characters")
	}

	db, err := sql.Open("pgx", *dsn)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("ping: %v", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("bcrypt: %v", err)
	}

	var id string
	err = db.QueryRowContext(ctx, "SELECT id FROM users WHERE username = $1", *username).Scan(&id)

	if err == nil {
		if !*force {
			fmt.Printf("user %s already exists (id=%s). Use -force to reset password.\n", *username, id)
			fmt.Printf("     username: %s | password: %s | role: %s\n", *username, *password, *role)
			return
		}
		_, uerr := db.ExecContext(ctx,
			"UPDATE users SET password_hash=$1, full_name=$2, email=$3, role=$4, is_active=TRUE, updated_at=NOW() WHERE id=$5",
			string(hash), *fullName, *email, *role, id,
		)
		if uerr != nil {
			log.Fatalf("update: %v", uerr)
		}
		fmt.Printf("OK — user %s password reset (id=%s role=%s)\n", *username, id, *role)
		fmt.Printf("     username: %s | password: %s\n", *username, *password)
		return
	}
	if err != sql.ErrNoRows {
		log.Fatalf("query: %v", err)
	}

	_, err = db.ExecContext(ctx,
		`INSERT INTO users (id, username, password_hash, full_name, email, role, is_active, created_at, updated_at)
		 VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, TRUE, NOW(), NOW())`,
		*username, string(hash), *fullName, *email, *role,
	)
	if err != nil {
		log.Fatalf("insert: %v", err)
	}
	fmt.Printf("OK — user %s created (role=%s)\n", *username, *role)
	fmt.Printf("     username: %s | password: %s\n", *username, *password)
}

func validateRole(r string) error {
	switch r {
	case "admin", "chief_accountant", "accountant", "viewer":
		return nil
	default:
		return fmt.Errorf("invalid role %q — must be admin | chief_accountant | accountant | viewer", r)
	}
}
