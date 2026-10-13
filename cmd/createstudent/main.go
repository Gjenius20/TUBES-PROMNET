// Perintah ini membuat akun STUDENT baru.
//
//	go run ./cmd/createstudent -name "Nama Siswa" -email siswa@example.com -nim 123456 -password "password123"
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
	name := flag.String("name", "", "nama siswa (wajib)")
	email := flag.String("email", "", "email siswa (wajib)")
	nim := flag.String("nim", "", "NIM siswa (wajib)")
	password := flag.String("password", "", "password siswa (8-72 karakter)")
	flag.Parse()

	*email = strings.ToLower(strings.TrimSpace(*email))
	*nim = strings.TrimSpace(*nim)
	if *name == "" || *email == "" || *nim == "" || len(*password) < 8 || len(*password) > 72 {
		log.Fatal("usage: createstudent -name <name> -email <email> -nim <nim> -password <8-72 chars>")
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
		Role:         models.RoleStudent,
	}
	if err := users.Create(context.Background(), user); err != nil {
		if errors.Is(err, repositories.ErrDuplicate) {
			log.Fatalf("a user with email %s or NIM %s already exists", *email, *nim)
		}
		log.Fatalf("create user: %v", err)
	}

	log.Printf("✅ Akun siswa '%s' (%s) berhasil dibuat!", user.Name, user.Email)
}
