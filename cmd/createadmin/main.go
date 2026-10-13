// Perintah ini membuat akun ADMIN (registrasi publik hanya membuat STUDENT).
//
//	go run ./cmd/createadmin -name "Admin" -email admin@example.com -nim 123456 -password "min-8-karakter"
package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"strings"

	"autograder/internal/config"
	"autograder/internal/models"
	"autograder/internal/repositories"
	"autograder/internal/utils"

	"github.com/joho/godotenv"
)

func main() {
	name := flag.String("name", "Admin", "nama admin")
	email := flag.String("email", "", "email admin")
	nim := flag.String("nim", "", "NIM/NISN admin (wajib)")
	password := flag.String("password", "", "password admin (8-72 karakter)")
	flag.Parse()

	*email = strings.ToLower(strings.TrimSpace(*email))
	*nim = strings.TrimSpace(*nim)
	if *email == "" || *nim == "" || len(*password) < 8 || len(*password) > 72 {
		log.Fatal("usage: createadmin -email <email> -nim <nim> -password <8-72 chars> [-name <name>]")
	}

	_ = godotenv.Load()
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN is not set")
	}

	db, err := config.ConnectDatabase(dsn)
	if err != nil {
		log.Fatalf("%v", err)
	}
	users := repositories.NewUserRepository(db)

	hash, err := utils.HashPassword(*password)
	if err != nil {
		log.Fatalf("hash password: %v", err)
	}

	user := &models.User{
		Name:         strings.TrimSpace(*name),
		Email:        *email,
		NIM:          *nim,
		PasswordHash: hash,
		Role:         models.RoleAdmin,
	}
	if err := users.Create(context.Background(), user); err != nil {
		if errors.Is(err, repositories.ErrDuplicate) {
			log.Fatalf("a user with email %s or NIM %s already exists", *email, *nim)
		}
		log.Fatalf("create admin: %v", err)
	}
	log.Printf("admin %s (%s) created (id=%d)", user.Email, user.NIM, user.ID)
}
