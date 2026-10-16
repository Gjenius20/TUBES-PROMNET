package main

import (
	"errors"
	"log"
	"math/rand"
	"time"

	"autograder/internal/config"
	"autograder/internal/models"
	"autograder/internal/utils"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	// 1. Load Environment
	_ = godotenv.Load()
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Gagal load config: %v", err)
	}

	// 2. Koneksi Database
	db, err := config.ConnectDatabase(cfg.DBDSN)
	if err != nil {
		log.Fatalf("❌ Gagal koneksi database: %v", err)
	}
	log.Println("✅ Database terkoneksi")

	// 3. Seed Data
	runSeed(db)

	log.Println("🎉 Seeding selesai dengan sukses!")
}

func runSeed(db *gorm.DB) {
	rand.Seed(time.Now().UnixNano())

	// --- Helper: Hash Password ---
	passwordHash, _ := utils.HashPassword("password123") // Default password untuk semua dummy

	// --- 1. Buat 5 Mahasiswa (Role STUDENT) ---
	students := []*models.User{
		{Name: "Budi Santoso", Email: "budi@student.com", NIM: "20240001", PasswordHash: passwordHash, Role: models.RoleStudent},
		{Name: "Siti Aminah", Email: "siti@student.com", NIM: "20240002", PasswordHash: passwordHash, Role: models.RoleStudent},
		{Name: "Joko Widodo", Email: "joko@student.com", NIM: "20240003", PasswordHash: passwordHash, Role: models.RoleStudent},
		{Name: "Dewi Lestari", Email: "dewi@student.com", NIM: "20240004", PasswordHash: passwordHash, Role: models.RoleStudent},
		{Name: "Ahmad Fauzi", Email: "ahmad@student.com", NIM: "20240005", PasswordHash: passwordHash, Role: models.RoleStudent},
	}

	for _, s := range students {
		// Gunakan FirstOrCreate agar idempotent (bisa dijalankan berulang tanpa error duplicate)
		if err := db.Where(models.User{Email: s.Email}).FirstOrCreate(s).Error; err != nil {
			log.Printf("⚠️ Gagal seed user %s: %v", s.Email, err)
		} else {
			log.Printf("👤 User created/found: %s (%s) ID:%d", s.Name, s.NIM, s.ID)
		}
	}

	// Ambil user admin yang sudah ada (ID=1 biasanya) sebagai referensi
	var admin models.User
	err := db.First(&admin, "role = ?", models.RoleAdmin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("👑 Admin tidak ditemukan, membuat admin default...")
			adminHash, _ := utils.HashPassword("admin123")
			admin = models.User{
				Name:         "Admin Default",
				Email:        "admin@cgrader.local",
				NIM:          "ADMIN001",
				PasswordHash: adminHash,
				Role:         models.RoleAdmin,
			}
			if err := db.Create(&admin).Error; err != nil {
				log.Fatalf("❌ Gagal membuat admin default: %v", err)
			}
			log.Printf("👑 Admin default dibuat: %s (ID: %d)", admin.Name, admin.ID)
		} else {
			log.Fatalf("❌ Error query admin: %v", err)
		}
	} else {
		log.Printf("👑 Admin ditemukan: %s (ID: %d)", admin.Name, admin.ID)
	}

	// --- 2. Buat Course & Quiz ---
	courses := []models.Course{
		{
			Title:       "Dasar Pemrograman C",
			Description: "Materi pengenalan variabel, tipe data, percabangan, dan perulangan.",
		},
		{
			Title:       "Struktur Data Dasar",
			Description: "Array, Pointer, Struct, dan Linked List sederhana.",
		},
	}

	for i := range courses {
		if err := db.Create(&courses[i]).Error; err != nil {
			log.Printf("⚠️ Gagal seed course %s: %v", courses[i].Title, err)
		} else {
			log.Printf("📚 Course created: %s (ID: %d)", courses[i].Title, courses[i].ID)
		}
	}

	// Quiz untuk Course 1 (Dasar Pemrograman C)
	quizzes := []models.Quiz{
		{
			CourseID:    courses[0].ID,
			Title:       "Kuis 1: Hello World & Variabel",
			Description: "Buat program C yang mencetak 'Hello, World!' dan deklarasikan variabel integer `umur` bernilai 20, lalu cetak nilainya.",
			Constraints: "Gunakan printf. Tidak boleh menggunakan scanf.",
		},
		{
			CourseID:    courses[0].ID,
			Title:       "Kuis 2: Percabangan If-Else",
			Description: "Baca integer `n` dari stdin. Cetak 'Positif' jika n > 0, 'Negatif' jika n < 0, 'Nol' jika n == 0.",
			Constraints: "Gunakan if-else if-else.",
		},
		{
			CourseID:    courses[1].ID,
			Title:       "Kuis 3: Perulangan & Array",
			Description: "Baca 5 integer ke array. Cetak jumlah total dan rata-rata (float 2 desimal).",
			Constraints: "Gunakan for loop. Output rata-rata dengan format %.2f.",
		},
	}

	for i := range quizzes {
		if err := db.Create(&quizzes[i]).Error; err != nil {
			log.Printf("⚠️ Gagal seed quiz %s: %v", quizzes[i].Title, err)
		} else {
			log.Printf("❓ Quiz created: %s (ID: %d)", quizzes[i].Title, quizzes[i].ID)
		}
	}

	// --- 3. Buat TestCase untuk setiap Quiz (Minimal 1 agar submission bisa dinilai) ---
	testCases := []models.TestCase{
		// Quiz 1
		{QuizID: quizzes[0].ID, Stdin: "", ExpectedOutput: "Hello, World!\nUmur: 20", IsHidden: false},
		// Quiz 2
		{QuizID: quizzes[1].ID, Stdin: "10", ExpectedOutput: "Positif", IsHidden: false},
		{QuizID: quizzes[1].ID, Stdin: "-5", ExpectedOutput: "Negatif", IsHidden: true},
		{QuizID: quizzes[1].ID, Stdin: "0", ExpectedOutput: "Nol", IsHidden: true},
		// Quiz 3
		{QuizID: quizzes[2].ID, Stdin: "10\n20\n30\n40\n50", ExpectedOutput: "Total: 150\nRata-rata: 30.00", IsHidden: false},
	}

	for _, tc := range testCases {
		if err := db.Create(&tc).Error; err != nil {
			log.Printf("⚠️ Gagal seed testcase untuk Quiz %d: %v", tc.QuizID, err)
		}
	}
	log.Printf("🧪 Testcases created: %d", len(testCases))

	// --- 4. Enroll Students ke Courses ---
	// Sesuai AGENTS.md: tidak semua siswa terdaftar di semua course
	statusApproved := models.EnrollmentApproved
	// Siswa 1 & 2 enroll di Course 1 saja
	// Siswa 3 & 4 enroll di Course 2 saja
	// Siswa 5 tidak terdaftar di course manapun
	enrollments := []struct {
		StudentID uint
		CourseID  uint
	}{
		{StudentID: students[0].ID, CourseID: courses[0].ID}, // Budi -> Course 1
		{StudentID: students[1].ID, CourseID: courses[0].ID}, // Siti -> Course 1
		{StudentID: students[2].ID, CourseID: courses[1].ID}, // Joko -> Course 2
		{StudentID: students[3].ID, CourseID: courses[1].ID}, // Dewi -> Course 2
		// Ahmad (students[4]) tidak terdaftar
	}
	for _, e := range enrollments {
		enroll := models.Enrollment{
			UserID:     e.StudentID,
			CourseID:   e.CourseID,
			Status:     statusApproved,
			EnrolledBy: &admin.ID, // Admin yang mendaftarkan
		}
		db.Where(models.Enrollment{UserID: e.StudentID, CourseID: e.CourseID}).FirstOrCreate(&enroll)
	}
	log.Println("📝 Enrollment created for students per AGENTS.md scenario.")

	// --- 5. Simulasi Submission Acak ---
	statuses := []models.SubmissionStatus{models.StatusPassed, models.StatusWrongAnswer}
	sourceCodes := map[uint]string{
		quizzes[0].ID: `#include <stdio.h>
int main() {
    int umur = 20;
    printf("Hello, World!\n");
    printf("Umur: %d", umur);
    return 0;
}`,
		quizzes[1].ID: `#include <stdio.h>
int main() {
    int n; scanf("%d", &n);
    if (n > 0) printf("Positif");
    else if (n < 0) printf("Negatif");
    else printf("Nol");
    return 0;
}`,
		quizzes[2].ID: `#include <stdio.h>
int main() {
    int arr[5], sum = 0;
    for(int i=0; i<5; i++) { scanf("%d", &arr[i]); sum += arr[i]; }
    printf("Total: %d\nRata-rata: %.2f", sum, sum/5.0);
    return 0;
}`,
	}

	submissionCount := 0
	for _, s := range students {
		for _, q := range quizzes {
			// Setiap siswa submit 1-2 kali per quiz
			for attempt := 0; attempt < rand.Intn(2)+1; attempt++ {
				status := statuses[rand.Intn(len(statuses))]
				score := 0
				passedCount := 0
				totalCount := 0

				// Hitung total testcase untuk quiz ini
				var totalCount64 int64
				db.Model(&models.TestCase{}).Where("quiz_id = ?", q.ID).Count(&totalCount64)
				totalCount = int(totalCount64)

				if status == models.StatusPassed {
					score = 100
					passedCount = totalCount
				} else {
					score = rand.Intn(80) // 0-79
					passedCount = rand.Intn(totalCount)
				}

				sub := models.Submission{
					UserID:      s.ID,
					QuizID:      q.ID,
					SourceCode:  sourceCodes[q.ID],
					Status:      status,
					Score:       score,
					PassedCount: passedCount,
					TotalCount:  totalCount,
					ResultsJson: "[]", // Valid JSON empty array untuk kolom JSON MySQL
				}

				if err := db.Create(&sub).Error; err != nil {
					log.Printf("⚠️ Gagal seed submission: %v", err)
				} else {
					submissionCount++
				}
			}
		}
	}

	log.Printf("📤 Submissions created: %d", submissionCount)
}