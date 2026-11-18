package main

import (
	"context"
	"log"
	"time"

	"github.com/geniussociety/backend/internal/config"
	"github.com/geniussociety/backend/internal/model"
	"github.com/geniussociety/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type SeedUser struct {
	Name     string
	Username string
	Email    string
	Password string
	Role     model.Role
	IsWalas  bool
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load failed: %v", err)
	}

	db, err := config.ConnectMongo(cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		log.Fatalf("mongo connect failed: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	teacherRepo := repository.NewTeacherRepository(db)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	seedUsers := []SeedUser{
		{
			Name:     "Super Administrator",
			Username: "superadmin",
			Email:    "superadmin@genius.ac.id",
			Password: "Password123!",
			Role:     model.RoleSuperAdmin,
		},
		{
			Name:     "Administrator Sekolah",
			Username: "admin",
			Email:    "admin@genius.ac.id",
			Password: "Password123!",
			Role:     model.RoleAdmin,
		},
		{
			Name:     "Budi Rahardjo, M.Kom",
			Username: "guru1",
			Email:    "budi@genius.ac.id",
			Password: "Password123!",
			Role:     model.RoleTeacher,
			IsWalas:  true,
		},
		{
			Name:     "Siti Aminah, S.Pd",
			Username: "guru2",
			Email:    "siti@genius.ac.id",
			Password: "Password123!",
			Role:     model.RoleTeacher,
			IsWalas:  false,
		},
		{
			Name:     "Drs. Hendro Wibowo",
			Username: "kurikulum",
			Email:    "kurikulum@genius.ac.id",
			Password: "Password123!",
			Role:     model.RoleCurriculum,
		},
		{
			Name:     "Dr. H. Bambang Soeprapto, M.Pd",
			Username: "kepsek",
			Email:    "kepsek@genius.ac.id",
			Password: "Password123!",
			Role:     model.RolePrincipal,
		},
		{
			Name:     "Ahmad Fauzi",
			Username: "murid1",
			Email:    "murid1@genius.ac.id",
			Password: "Password123!",
			Role:     model.RoleStudent,
		},
		{
			Name:     "Putri Anggraini",
			Username: "murid2",
			Email:    "murid2@genius.ac.id",
			Password: "Password123!",
			Role:     model.RoleStudent,
		},
	}

// Create or fetch demo school
schoolRepo := repository.NewSchoolRepository(db)
	demoSchool, err := schoolRepo.FindOrCreateByName(ctx, "Demo School", "Jl. Pendidikan No.1, Jakarta")
	if err != nil {
		log.Fatalf("failed to get or create demo school: %v", err)
	}
	demoSchoolID := demoSchool.ID
	log.Printf("Demo school ID: %s", demoSchoolID.Hex())

	// Inside cmd/seed/main.go – after creating/finding the demo school
	// ----- Academic Year seed (demo) -----
	demoAYName := "2026/2027"
	ayRepo := repository.NewAcademicYearRepository(db)
	existingAY, _ := ayRepo.FindByName(context.TODO(), demoSchool.ID, demoAYName)
	if existingAY == nil {
		newAY := &model.AcademicYear{
			SchoolID:  demoSchool.ID,
			Name:      demoAYName,
			Semester:  1,
			StartDate: time.Date(2026, time.July, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2027, time.June, 30, 23, 59, 59, 0, time.UTC),
			IsActive:  true,
		}
		if err := ayRepo.Create(context.TODO(), newAY); err != nil {
			log.Printf("❗️ failed to create demo academic year: %v", err)
		} else {
			log.Printf("✅ demo academic year created: %s (ID %s)", demoAYName, newAY.ID.Hex())
		}
	} else {
		log.Printf("✅ demo academic year already exists: %s (ID %s)", demoAYName, existingAY.ID.Hex())
	}


	// Early user processing loop removed; handled by later comprehensive loop.



		// ---------------------------------------------------
		// Duplicate school creation and migration removed; handled earlier in script.


		// ---------------------------------------------------
		// 2️⃣ Insert / update users (including admin)
		for _, su := range seedUsers {
			existing, _ := userRepo.FindByUsername(ctx, su.Username)

			// ----- CREATE new user -----
			if existing == nil {
				// password hashing for new user
				hash, err := bcrypt.GenerateFromPassword([]byte(su.Password), bcrypt.DefaultCost)
				if err != nil {
					log.Fatalf("bcrypt error: %v", err)
				}
				u := &model.User{
					Name:         su.Name,
					Username:     su.Username,
					Email:        su.Email,
					PasswordHash: string(hash),
					Role:         su.Role,
					IsActive:     true,
				}

				// Assign demo school ID to all school-scoped users
				if su.Role != model.RoleSuperAdmin {
					u.SchoolID = demoSchoolID
				} else {
					u.SchoolID = demoSchoolID
				}

				if err := userRepo.Create(ctx, u); err != nil {
					log.Printf("failed to insert user %s: %v", su.Username, err)
					continue
				}
				// teacher creation for new user
				if su.Role == model.RoleTeacher {
					_ = teacherRepo.Create(ctx, &model.Teacher{UserID: u.ID, NIP: "19850101" + u.ID.Hex()[:6], IsWalas: su.IsWalas})
				}
				continue
			}

			// ----- UPDATE existing user when they lack a SchoolID -----
			if existing.SchoolID.IsZero() {
				// Set the proper SchoolID via an update call
				if err := userRepo.UpdateSchoolID(ctx, existing.ID, demoSchoolID); err != nil {
					log.Printf("failed to set schoolId for user %s: %v", su.Username, err)
				} else {
					log.Printf("user %s patched with schoolId %s", su.Username, demoSchoolID.Hex())
				}
			}
		}

		// Seed 3 Mata Pelajaran Dasar jika kosong
		subjectRepo := repository.NewSubjectRepository(db)
		initialSubjects := []model.Subject{
			{Code: "MTK-WJB", Name: "Matematika Wajib", Category: "Umum", GradeLevel: "Semua Tingkat", KKM: 75, Status: "active", Description: "Penguasaan logika matematika, kalkulus, dan aljabar"},
			{Code: "IND-WJB", Name: "Bahasa Indonesia", Category: "Umum", GradeLevel: "Semua Tingkat", KKM: 75, Status: "active", Description: "Literasi teks analitis, eksposisi, dan karya ilmiah"},
			{Code: "ENG-WJB", Name: "Bahasa Inggris", Category: "Umum", GradeLevel: "Semua Tingkat", KKM: 75, Status: "active", Description: "Komunikasi aktif bahasa Inggris lisan dan tulisan"},
		}
		for _, is := range initialSubjects {
			if ex, _ := subjectRepo.FindByCode(ctx, is.Code); ex == nil {
				_ = subjectRepo.Create(ctx, &is)
				log.Printf("Seeded subject: %s (%s)", is.Name, is.Code)
			}
		}

		log.Println("Seeding development completed successfully!")
	}

