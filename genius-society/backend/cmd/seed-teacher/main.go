package main

import (
	"context"
	"log"
	"time"

	"github.com/geniussociety/backend/internal/config"
	"github.com/geniussociety/backend/internal/model"
	"github.com/geniussociety/backend/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := config.ConnectMongo(cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		log.Fatalf("failed to connect mongo: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	userRepo := repository.NewUserRepository(db)
	teacherRepo := repository.NewTeacherRepository(db)
	schoolRepo := repository.NewSchoolRepository(db)
	ayRepo := repository.NewAcademicYearRepository(db)
	classRepo := repository.NewClassRepository(db)
	subjectRepo := repository.NewSubjectRepository(db)

	teacherClassesColl := db.Collection("teacher_classes")
	teacherSubjectsColl := db.Collection("teacher_subjects")
	teachersColl := db.Collection("teachers")
	classesColl := db.Collection("classes")

	// 1. Get or create demo school
	demoSchool, err := schoolRepo.FindOrCreateByName(ctx, "Demo School", "Jl. Pendidikan No.1, Jakarta")
	if err != nil {
		log.Fatalf("failed to get demo school: %v", err)
	}
	log.Printf("Demo school ID: %s", demoSchool.ID.Hex())

	// 2. Get or create active academic year
	ay, err := ayRepo.FindActiveBySchool(ctx, demoSchool.ID)
	if err != nil || ay == nil {
		demoAYName := "2026/2027"
		existingAY, _ := ayRepo.FindByName(ctx, demoSchool.ID, demoAYName)
		if existingAY != nil {
			ay = existingAY
		} else {
			ay = &model.AcademicYear{
				SchoolID:  demoSchool.ID,
				Name:      demoAYName,
				Semester:  1,
				StartDate: time.Date(2026, time.July, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2027, time.June, 30, 23, 59, 59, 0, time.UTC),
				IsActive:  true,
				Status:    "active",
			}
			if err := ayRepo.Create(ctx, ay); err != nil {
				log.Fatalf("failed to create academic year: %v", err)
			}
		}
	}
	log.Printf("Active academic year: %s (ID %s)", ay.Name, ay.ID.Hex())

	// 3. Find guru1 and guru2 users
	guru1User, err := userRepo.FindByUsername(ctx, "guru1")
	if err != nil || guru1User == nil {
		log.Fatalf("user guru1 not found in database; please run cmd/seed first")
	}
	guru2User, err := userRepo.FindByUsername(ctx, "guru2")
	if err != nil || guru2User == nil {
		log.Fatalf("user guru2 not found in database; please run cmd/seed first")
	}

	// Ensure users have demo school ID
	if guru1User.SchoolID.IsZero() {
		_ = userRepo.UpdateSchoolID(ctx, guru1User.ID, demoSchool.ID)
	}
	if guru2User.SchoolID.IsZero() {
		_ = userRepo.UpdateSchoolID(ctx, guru2User.ID, demoSchool.ID)
	}

	// 4. Find or create teacher profiles
	t1, err := teacherRepo.FindByUserID(ctx, guru1User.ID)
	if err != nil || t1 == nil {
		t1 = &model.Teacher{
			UserID:    guru1User.ID,
			NIP:       "19850101" + guru1User.ID.Hex()[:6],
			IsWalas:   true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := teacherRepo.Create(ctx, t1); err != nil {
			log.Fatalf("failed to create teacher guru1: %v", err)
		}
	}

	t2, err := teacherRepo.FindByUserID(ctx, guru2User.ID)
	if err != nil || t2 == nil {
		t2 = &model.Teacher{
			UserID:    guru2User.ID,
			NIP:       "19870202" + guru2User.ID.Hex()[:6],
			IsWalas:   false,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := teacherRepo.Create(ctx, t2); err != nil {
			log.Fatalf("failed to create teacher guru2: %v", err)
		}
	}

	// 5. Find or create classes X-IPA-1 and X-IPA-2
	classNames := []string{"X-IPA-1", "X-IPA-2"}
	classes := make(map[string]*model.Class)
	for _, name := range classNames {
		filter := bson.M{
			"schoolId":       demoSchool.ID,
			"academicYearId": ay.ID,
			"name":           name,
		}
		var cl model.Class
		err := classesColl.FindOne(ctx, filter).Decode(&cl)
		if err == nil {
			classes[name] = &cl
		} else {
			newClass := &model.Class{
				SchoolID:       demoSchool.ID,
				Code:           name,
				Name:           name,
				EducationLevel: "SMA",
				AcademicYearID: ay.ID,
				Capacity:       36,
				Status:         "active",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}
			if err := classRepo.Create(ctx, newClass); err != nil {
				log.Fatalf("failed to create class %s: %v", name, err)
			}
			classes[name] = newClass
			log.Printf("Created class: %s (ID %s)", name, newClass.ID.Hex())
		}
	}

	// 6. Configure Walas: guru1 is walas of X-IPA-1, guru2 is non-walas
	class1 := classes["X-IPA-1"]
	class2 := classes["X-IPA-2"]

	// Set homeroom on class1
	_, _ = classesColl.UpdateOne(ctx,
		bson.M{"_id": class1.ID},
		bson.M{"$set": bson.M{"homeroomTeacherId": t1.ID, "updatedAt": time.Now()}},
	)

	// Set walas on teachers
	_, _ = teachersColl.UpdateOne(ctx,
		bson.M{"_id": t1.ID},
		bson.M{"$set": bson.M{
			"isWalas":       true,
			"walasClassIds": []bson.ObjectID{class1.ID},
			"updatedAt":     time.Now(),
		}},
	)
	_, _ = teachersColl.UpdateOne(ctx,
		bson.M{"_id": t2.ID},
		bson.M{"$set": bson.M{
			"isWalas":       false,
			"walasClassIds": []bson.ObjectID{},
			"updatedAt":     time.Now(),
		}},
	)
	log.Printf("Configured Walas: guru1 -> %s (Walas: true), guru2 -> (Walas: false)", class1.Name)

	// 7. Find subjects (MTK-WJB, IND-WJB)
	subMtk, _ := subjectRepo.FindByCode(ctx, "MTK-WJB")
	subInd, _ := subjectRepo.FindByCode(ctx, "IND-WJB")

	if subMtk == nil {
		subMtk = &model.Subject{
			Code:        "MTK-WJB",
			Name:        "Matematika Wajib",
			Category:    "Umum",
			GradeLevel:  "Semua Tingkat",
			KKM:         75,
			Status:      "active",
			Description: "Penguasaan logika matematika",
		}
		_ = subjectRepo.Create(ctx, subMtk)
	}
	if subInd == nil {
		subInd = &model.Subject{
			Code:        "IND-WJB",
			Name:        "Bahasa Indonesia",
			Category:    "Umum",
			GradeLevel:  "Semua Tingkat",
			KKM:         75,
			Status:      "active",
			Description: "Literasi teks analitis",
		}
		_ = subjectRepo.Create(ctx, subInd)
	}

	// 8. Idempotently assign teacher_subjects
	assignSubject := func(teacherID, subjectID bson.ObjectID) {
		filter := bson.M{
			"teacherId": teacherID,
			"subjectId": subjectID,
		}
		update := bson.M{
			"$setOnInsert": bson.M{
				"_id":       bson.NewObjectID(),
				"teacherId": teacherID,
				"subjectId": subjectID,
				"schoolId":  demoSchool.ID,
				"createdAt": time.Now(),
				"updatedAt": time.Now(),
			},
		}
		opts := options.UpdateOne().SetUpsert(true)
		_, err := teacherSubjectsColl.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			log.Printf("failed to assign teacher_subject (%s, %s): %v", teacherID.Hex(), subjectID.Hex(), err)
		}
	}
	assignSubject(t1.ID, subMtk.ID)
	assignSubject(t2.ID, subInd.ID)
	log.Println("Assigned teacher_subjects: guru1 -> MTK-WJB, guru2 -> IND-WJB")

	// 9. Idempotently assign teacher_classes
	assignClass := func(teacherID, classID, subjectID bson.ObjectID) {
		filter := bson.M{
			"teacherId":      teacherID,
			"classId":        classID,
			"subjectId":      subjectID,
			"academicYearId": ay.ID,
		}
		update := bson.M{
			"$setOnInsert": bson.M{
				"_id":            bson.NewObjectID(),
				"teacherId":      teacherID,
				"classId":        classID,
				"subjectId":      subjectID,
				"academicYearId": ay.ID,
				"schoolId":       demoSchool.ID,
				"createdAt":      time.Now(),
				"updatedAt":      time.Now(),
			},
		}
		opts := options.UpdateOne().SetUpsert(true)
		_, err := teacherClassesColl.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			log.Printf("failed to assign teacher_class: %v", err)
		}
	}

	// guru1 teaches MTK in class1 (X-IPA-1) and class2 (X-IPA-2)
	assignClass(t1.ID, class1.ID, subMtk.ID)
	assignClass(t1.ID, class2.ID, subMtk.ID)

	// guru2 teaches IND in class1 (X-IPA-1) and class2 (X-IPA-2)
	assignClass(t2.ID, class1.ID, subInd.ID)
	assignClass(t2.ID, class2.ID, subInd.ID)

	log.Println("Assigned teacher_classes: guru1 teaches MTK in X-IPA-1 & X-IPA-2; guru2 teaches IND in X-IPA-1 & X-IPA-2")
	log.Println("✅ Teacher seeding completed successfully!")
}
