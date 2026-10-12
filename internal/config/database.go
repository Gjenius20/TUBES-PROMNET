package config

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ConnectDatabase membuka koneksi GORM ke MySQL dengan retry, karena container
// MySQL bisa memerlukan beberapa detik sebelum siap menerima koneksi.
func ConnectDatabase(dsn string) (*gorm.DB, error) {
	const maxAttempts = 10

	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		db, err := openAndPing(dsn)
		if err == nil {
			return db, nil
		}
		lastErr = err
		log.Printf("database not ready (attempt %d/%d): %v", attempt, maxAttempts, lastErr)
		time.Sleep(3 * time.Second)
	}
	return nil, fmt.Errorf("could not connect to database: %w", lastErr)
}

func openAndPing(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// Menerjemahkan error driver (mis. duplicate key) menjadi gorm.ErrDuplicatedKey.
		TranslateError: true,
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	return db, nil
}
