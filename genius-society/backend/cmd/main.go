package main

import (
	"fmt"
	"log"

	"github.com/geniussociety/backend/internal/config"
	"github.com/geniussociety/backend/internal/repository"
	"github.com/geniussociety/backend/routes"
	"github.com/geniussociety/backend/internal/service"
	"github.com/geniussociety/backend/internal/storage"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	db, err := config.ConnectMongo(cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		log.Fatalf("mongodb connection error: %v", err)
	}
	log.Println("Connected to MongoDB:", cfg.MongoDB)

	// Repositories
	userRepo := repository.NewUserRepository(db)
	teacherRepo := repository.NewTeacherRepository(db)
	studentRepo := repository.NewStudentRepository(db)
	classRepo := repository.NewClassRepository(db)
	subjectRepo := repository.NewSubjectRepository(db)

	// Services
	authSvc := service.NewAuthService(userRepo, teacherRepo, cfg)
	userSvc := service.NewUserService(userRepo, teacherRepo, studentRepo)
	adminDashboardSvc := service.NewAdminDashboardService(db)
	studentManagementSvc := service.NewStudentManagementService(studentRepo, userRepo, classRepo, db)
	teacherManagementSvc := service.NewTeacherManagementService(teacherRepo, userRepo, classRepo, db)
	subjectManagementSvc := service.NewSubjectManagementService(subjectRepo, teacherRepo, db)

	// Storage (B2) — optional: log warning if not configured
	b2, err := storage.NewB2Storage(cfg)
	if err != nil {
		log.Printf("WARNING: B2 storage not configured: %v", err)
		b2 = nil
	} else {
		log.Println("B2 storage initialized:", cfg.B2BucketName)
	}
	_ = b2 // will be injected into future handlers

	// Router
	r := routes.Setup(cfg, authSvc, userSvc, adminDashboardSvc, studentManagementSvc, teacherManagementSvc, subjectManagementSvc)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on %s (env: %s)", addr, cfg.AppEnv)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
