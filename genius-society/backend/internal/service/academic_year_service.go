package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/geniussociety/backend/internal/dto"
	"github.com/geniussociety/backend/internal/model"
	"github.com/geniussociety/backend/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AcademicYearService struct {
	ayRepo         *repository.AcademicYearRepository
	studentRepo    *repository.StudentRepository
	classRepo      *repository.ClassRepository
	userRepo       *repository.UserRepository
	enrollmentRepo *repository.EnrollmentRepository
	decisionRepo   *repository.AcademicDecisionRepository
	jobRepo        *repository.RolloverJobRepository
	auditRepo      *repository.AuditLogRepository
	db             *mongo.Database
}

func NewAcademicYearService(
	ayRepo *repository.AcademicYearRepository,
	studentRepo *repository.StudentRepository,
	classRepo *repository.ClassRepository,
	userRepo *repository.UserRepository,
	enrollmentRepo *repository.EnrollmentRepository,
	decisionRepo *repository.AcademicDecisionRepository,
	jobRepo *repository.RolloverJobRepository,
	auditRepo *repository.AuditLogRepository,
	db *mongo.Database,
) *AcademicYearService {
	return &AcademicYearService{
		ayRepo:         ayRepo,
		studentRepo:    studentRepo,
		classRepo:      classRepo,
		userRepo:       userRepo,
		enrollmentRepo: enrollmentRepo,
		decisionRepo:   decisionRepo,
		jobRepo:        jobRepo,
		auditRepo:      auditRepo,
		db:             db,
	}
}

// ---------- LIST ----------
func (s *AcademicYearService) List(ctx context.Context, schoolID bson.ObjectID) ([]dto.AcademicYearResponse, error) {
	years, err := s.ayRepo.FindAllBySchool(ctx, schoolID)
	if err != nil {
		return nil, err
	}
	resp := make([]dto.AcademicYearResponse, 0, len(years))
	for _, y := range years {
		status := y.Status
		if status == "" {
			if y.IsActive {
				status = "active"
			} else {
				status = "inactive"
			}
		}
		resp = append(resp, dto.AcademicYearResponse{
			ID:        y.ID.Hex(),
			Name:      y.Name,
			Semester:  y.Semester,
			StartDate: y.StartDate,
			EndDate:   y.EndDate,
			IsActive:  y.IsActive,
			Status:    status,
			CreatedAt: y.CreatedAt,
			UpdatedAt: y.UpdatedAt,
		})
	}
	return resp, nil
}

// ---------- GET ----------
func (s *AcademicYearService) GetOne(ctx context.Context, schoolID, id bson.ObjectID) (*dto.AcademicYearResponse, error) {
	ay, err := s.ayRepo.FindByID(ctx, schoolID, id)
	if err != nil {
		return nil, err
	}
	status := ay.Status
	if status == "" {
		if ay.IsActive {
			status = "active"
		} else {
			status = "inactive"
		}
	}
	return &dto.AcademicYearResponse{
		ID:        ay.ID.Hex(),
		Name:      ay.Name,
		Semester:  ay.Semester,
		StartDate: ay.StartDate,
		EndDate:   ay.EndDate,
		IsActive:  ay.IsActive,
		Status:    status,
		CreatedAt: ay.CreatedAt,
		UpdatedAt: ay.UpdatedAt,
	}, nil
}

// ---------- CREATE ----------
func (s *AcademicYearService) Create(ctx context.Context, schoolID, userID bson.ObjectID, req dto.CreateAcademicYearDTO) (*dto.AcademicYearResponse, error) {
	// Duplicate-name guard
	if existing, _ := s.ayRepo.FindByName(ctx, schoolID, req.Name); existing != nil {
		return nil, errors.New("tahun ajaran dengan nama tersebut sudah ada")
	}
	if req.EndDate.Before(req.StartDate) {
		return nil, errors.New("tanggal akhir harus setelah tanggal mulai")
	}
	ay := &model.AcademicYear{
		SchoolID:  schoolID,
		Name:      req.Name,
		Semester:  req.Semester,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		IsActive:  false,
		Status:    "inactive",
	}
	if err := s.ayRepo.Create(ctx, ay); err != nil {
		return nil, err
	}
	if s.auditRepo != nil {
		_ = s.auditRepo.Log(ctx, schoolID, userID, "academic_year.create", "academic_year", ay.ID.Hex(), bson.M{"name": ay.Name})
	}
	return &dto.AcademicYearResponse{
		ID:        ay.ID.Hex(),
		Name:      ay.Name,
		Semester:  ay.Semester,
		StartDate: ay.StartDate,
		EndDate:   ay.EndDate,
		IsActive:  ay.IsActive,
		Status:    ay.Status,
		CreatedAt: ay.CreatedAt,
		UpdatedAt: ay.UpdatedAt,
	}, nil
}

// ---------- UPDATE ----------
func (s *AcademicYearService) Update(ctx context.Context, schoolID, userID, id bson.ObjectID, req dto.UpdateAcademicYearDTO) (*dto.AcademicYearResponse, error) {
	existing, err := s.ayRepo.FindByID(ctx, schoolID, id)
	if err != nil {
		return nil, errors.New("tahun ajaran tidak ditemukan")
	}
	if existing.Status == "archived" {
		return nil, errors.New("tahun ajaran yang diarsipkan tidak dapat diubah")
	}

	update := bson.M{}
	if req.Name != nil {
		if *req.Name == "" {
			return nil, errors.New("nama tidak boleh kosong")
		}
		if dup, _ := s.ayRepo.FindByName(ctx, schoolID, *req.Name); dup != nil && dup.ID != id {
			return nil, errors.New("nama tahun ajaran sudah dipakai")
		}
		update["name"] = *req.Name
	}
	if req.Semester != nil {
		if *req.Semester != 1 && *req.Semester != 2 {
			return nil, errors.New("semester harus 1 atau 2")
		}
		update["semester"] = *req.Semester
	}
	if req.StartDate != nil {
		update["startDate"] = *req.StartDate
	}
	if req.EndDate != nil {
		update["endDate"] = *req.EndDate
	}
	if sd, ok1 := update["startDate"]; ok1 {
		if ed, ok2 := update["endDate"]; ok2 && ed.(time.Time).Before(sd.(time.Time)) {
			return nil, errors.New("tanggal akhir harus setelah tanggal mulai")
		}
	}
	if len(update) == 0 {
		return nil, errors.New("tidak ada data yang diubah")
	}
	if err := s.ayRepo.Update(ctx, schoolID, id, update); err != nil {
		return nil, err
	}
	if s.auditRepo != nil {
		_ = s.auditRepo.Log(ctx, schoolID, userID, "academic_year.update", "academic_year", id.Hex(), update)
	}
	return s.GetOne(ctx, schoolID, id)
}

// ---------- ACTIVATE ----------
func (s *AcademicYearService) Activate(ctx context.Context, schoolID, userID, id bson.ObjectID) error {
	existing, err := s.ayRepo.FindByID(ctx, schoolID, id)
	if err != nil {
		return errors.New("tahun ajaran tidak ditemukan")
	}
	if existing.Status == "archived" {
		return errors.New("tahun ajaran yang diarsipkan tidak dapat diaktifkan kembali secara langsung")
	}
	if err := s.ayRepo.Activate(ctx, schoolID, id); err != nil {
		return err
	}
	if s.auditRepo != nil {
		_ = s.auditRepo.Log(ctx, schoolID, userID, "academic_year.activate", "academic_year", id.Hex(), bson.M{"name": existing.Name})
	}
	return nil
}

// ---------- ARCHIVE ----------
func (s *AcademicYearService) Archive(ctx context.Context, schoolID, userID, id bson.ObjectID) error {
	existing, err := s.ayRepo.FindByID(ctx, schoolID, id)
	if err != nil {
		return errors.New("tahun ajaran tidak ditemukan")
	}
	if existing.IsActive {
		return errors.New("tidak dapat mengarsipkan tahun ajaran yang sedang aktif; aktifkan tahun ajaran lain terlebih dahulu")
	}
	if err := s.ayRepo.Archive(ctx, schoolID, id); err != nil {
		return err
	}
	if s.auditRepo != nil {
		_ = s.auditRepo.Log(ctx, schoolID, userID, "academic_year.archive", "academic_year", id.Hex(), bson.M{"name": existing.Name})
	}
	return nil
}

// ---------- DELETE ----------
func (s *AcademicYearService) Delete(ctx context.Context, schoolID, userID, id bson.ObjectID) error {
	existing, err := s.ayRepo.FindByID(ctx, schoolID, id)
	if err != nil {
		return errors.New("tahun ajaran tidak ditemukan")
	}
	if existing.IsActive {
		return errors.New("tidak dapat menghapus tahun ajaran yang sedang aktif")
	}
	if s.classRepo != nil {
		cnt, _ := s.classRepo.CountByFilter(ctx, bson.M{"schoolId": schoolID, "academicYearId": id})
		if cnt > 0 {
			return errors.New("tahun ajaran ini sudah memiliki kelas terdaftar; gunakan fitur arsip dan bukan hapus")
		}
	}
	if err := s.ayRepo.Delete(ctx, schoolID, id); err != nil {
		return err
	}
	if s.auditRepo != nil {
		_ = s.auditRepo.Log(ctx, schoolID, userID, "academic_year.delete", "academic_year", id.Hex(), bson.M{"name": existing.Name})
	}
	return nil
}

// ============================================================================
// PROMOTION / ROLLOVER
// ============================================================================

// GetPromotionCandidates loads all active students belonging to classes in sourceYearID
func (s *AcademicYearService) GetPromotionCandidates(ctx context.Context, schoolID, sourceYearID bson.ObjectID) ([]dto.PromotionCandidateItem, error) {
	_, err := s.ayRepo.FindByID(ctx, schoolID, sourceYearID)
	if err != nil {
		return nil, errors.New("tahun ajaran asal tidak ditemukan")
	}

	// 1. Find classes for source academic year in this school
	classes, _, err := s.classRepo.List(ctx, bson.M{
		"schoolId":       schoolID,
		"academicYearId": sourceYearID,
		"status":         bson.M{"$ne": "archived"},
	}, 1, 1000)
	if err != nil {
		return nil, err
	}

	classMap := make(map[bson.ObjectID]*model.Class)
	classIDs := make([]bson.ObjectID, 0, len(classes))
	for _, c := range classes {
		classMap[c.ID] = c
		classIDs = append(classIDs, c.ID)
	}

	if len(classIDs) == 0 {
		return []dto.PromotionCandidateItem{}, nil
	}

	// 2. Find active students in these classes
	studentFilter := bson.M{
		"classId": bson.M{"$in": classIDs},
		"status":  "active",
	}
	students, _, err := s.studentRepo.List(ctx, studentFilter, 1, 5000)
	if err != nil {
		return nil, err
	}

	// 3. User map for names
	userIDs := make([]bson.ObjectID, 0)
	for _, st := range students {
		if !st.UserID.IsZero() {
			userIDs = append(userIDs, st.UserID)
		}
	}
	userMap := make(map[bson.ObjectID]*model.User)
	if len(userIDs) > 0 {
		cur, err := s.db.Collection("users").Find(ctx, bson.M{"_id": bson.M{"$in": userIDs}})
		if err == nil {
			var users []model.User
			if err := cur.All(ctx, &users); err == nil {
				for i := range users {
					userMap[users[i].ID] = &users[i]
				}
			}
		}
	}

	// 4. Calculate academic evaluation per student
	candidates := make([]dto.PromotionCandidateItem, 0, len(students))
	for _, st := range students {
		cls := classMap[st.ClassID]
		clsName := "-"
		eduLevel := "-"
		if cls != nil {
			clsName = cls.Name
			eduLevel = cls.EducationLevel
		}

		studentName := fmt.Sprintf("Siswa %s", st.NIS)
		if u, exists := userMap[st.UserID]; exists && u.Name != "" {
			studentName = u.Name
		}

		avgScore, missingGrades := s.evaluateStudentGrades(ctx, st.ID, st.ClassID, sourceYearID)

		eligibility := "eligible"
		recommended := "promoted"
		reasonCode := ""
		reason := ""

		if missingGrades < 0 {
			eligibility = "pending_review"
			recommended = "suspended"
			reasonCode = "no_grades"
			reason = "Belum ada nilai"
		} else if missingGrades > 0 {
			eligibility = "pending_review"
			recommended = "suspended"
			reasonCode = "grades_incomplete"
			reason = fmt.Sprintf("Nilai belum lengkap (%d komponen evaluasi belum terisi)", missingGrades)
		} else if avgScore < 60.0 && avgScore > 0 {
			eligibility = "pending_review"
			recommended = "retained"
			reasonCode = "academic_below_threshold"
			reason = fmt.Sprintf("Rata-rata nilai (%.1f) di bawah batas kriteria ketuntasan (KKM)", avgScore)
		}

		candidates = append(candidates, dto.PromotionCandidateItem{
			StudentID:          st.ID.Hex(),
			NIS:                st.NIS,
			Name:               studentName,
			SourceClassID:      st.ClassID.Hex(),
			SourceClassName:    clsName,
			EducationLevel:     eduLevel,
			EligibilityStatus:  eligibility,
			RecommendedAction:  recommended,
			ReasonCode:         reasonCode,
			Reason:             reason,
			AverageScore:       avgScore,
			MissingGradesCount: missingGrades,
		})
	}

	return candidates, nil
}

// PreviewPromotion validates decisions and returns a preview summary
func (s *AcademicYearService) PreviewPromotion(ctx context.Context, schoolID bson.ObjectID, req dto.PromotionPreviewRequest) (*dto.ProcessPreviewResponse, error) {
	sourceOID, err := bson.ObjectIDFromHex(req.SourceAcademicYearID)
	if err != nil {
		return nil, errors.New("id tahun ajaran asal tidak valid")
	}
	targetOID, err := bson.ObjectIDFromHex(req.TargetAcademicYearID)
	if err != nil {
		return nil, errors.New("id tahun ajaran tujuan tidak valid")
	}
	if sourceOID == targetOID {
		return nil, errors.New("tahun ajaran asal dan tujuan tidak boleh sama")
	}

	// Validate target academic year exists and belongs to school
	_, err = s.ayRepo.FindByID(ctx, schoolID, targetOID)
	if err != nil {
		return nil, errors.New("tahun ajaran tujuan tidak ditemukan dalam sekolah ini")
	}

	candidates, err := s.GetPromotionCandidates(ctx, schoolID, sourceOID)
	if err != nil {
		return nil, err
	}

	// Decision override map
	decMap := make(map[string]dto.PromotionDecisionInput)
	for _, d := range req.Decisions {
		decMap[d.StudentID] = d
	}

	var eligibleCount, pendingCount, promotedCount, retainedCount, suspendedCount, excludedCount int64
	var previewErrors []string

	for i := range candidates {
		c := &candidates[i]
		if d, ok := decMap[c.StudentID]; ok {
			c.RecommendedAction = d.Action
			c.TargetClassID = d.TargetClassID
			if d.Reason != "" {
				c.Reason = d.Reason
			}
			if d.ReasonCode != "" {
				c.ReasonCode = d.ReasonCode
			}
		}

		switch c.RecommendedAction {
		case "promoted":
			promotedCount++
			if c.TargetClassID == "" {
				previewErrors = append(previewErrors, fmt.Sprintf("Siswa %s (%s) dipilih naik kelas tetapi belum memiliki kelas tujuan", c.Name, c.NIS))
			}
		case "retained":
			retainedCount++
			if c.TargetClassID == "" {
				previewErrors = append(previewErrors, fmt.Sprintf("Siswa %s (%s) dipilih tinggal kelas tetapi belum memiliki kelas tujuan", c.Name, c.NIS))
			}
		case "suspended":
			suspendedCount++
			if c.Reason == "" {
				c.Reason = "Keputusan akademik ditangguhkan untuk verifikasi data lanjutan"
			}
		case "excluded":
			excludedCount++
		}

		if c.EligibilityStatus == "eligible" {
			eligibleCount++
		} else {
			pendingCount++
		}
	}

	return &dto.ProcessPreviewResponse{
		TotalStudents: int64(len(candidates)),
		Eligible:      eligibleCount,
		PendingReview: pendingCount,
		Promoted:      promotedCount,
		Retained:      retainedCount,
		Suspended:     suspendedCount,
		Excluded:      excludedCount,
		Errors:        previewErrors,
		Candidates:    candidates,
	}, nil
}

// ExecutePromotion executes the batch promotion with idempotency and audit logs
func (s *AcademicYearService) ExecutePromotion(ctx context.Context, schoolID, userID bson.ObjectID, req dto.PromotionExecutionRequest) (*dto.ProcessExecutionResponse, error) {
	if req.IdempotencyKey == "" {
		return nil, errors.New("idempotencyKey wajib disertakan")
	}

	// 1. Idempotency Check
	existingJob, err := s.jobRepo.FindByIdempotencyKey(ctx, schoolID, req.IdempotencyKey)
	if err == nil && existingJob != nil {
		return &dto.ProcessExecutionResponse{
			JobID: existingJob.ID.Hex(),
			Status: string(existingJob.Status),
			Summary: dto.RolloverJobSummaryDTO{
				TotalStudents: existingJob.Summary.TotalStudents,
				Promoted:      existingJob.Summary.Promoted,
				Retained:      existingJob.Summary.Retained,
				Graduated:     existingJob.Summary.Graduated,
				Suspended:     existingJob.Summary.Suspended,
				Excluded:      existingJob.Summary.Excluded,
				Failed:        existingJob.Summary.Failed,
			},
			ExecutedAt: existingJob.CreatedAt,
		}, nil
	}

	sourceOID, err := bson.ObjectIDFromHex(req.SourceAcademicYearID)
	if err != nil {
		return nil, errors.New("id tahun ajaran asal tidak valid")
	}
	targetOID, err := bson.ObjectIDFromHex(req.TargetAcademicYearID)
	if err != nil {
		return nil, errors.New("id tahun ajaran tujuan tidak valid")
	}
	if sourceOID == targetOID {
		return nil, errors.New("tahun ajaran asal dan tujuan tidak boleh sama")
	}

	// 2. Validate Target Academic Year
	targetYear, err := s.ayRepo.FindByID(ctx, schoolID, targetOID)
	if err != nil {
		return nil, errors.New("tahun ajaran tujuan tidak ditemukan dalam sekolah ini")
	}
	if targetYear.Status == "archived" {
		return nil, errors.New("tahun ajaran tujuan sudah diarsipkan")
	}

	// 3. Create RolloverJob in Processing state
	job := &model.RolloverJob{
		SchoolID:             schoolID,
		SourceAcademicYearID: sourceOID,
		TargetAcademicYearID: &targetOID,
		ProcessType:          model.ProcessTypePromotion,
		Status:               model.JobStatusProcessing,
		ExecutedBy:           userID,
		IdempotencyKey:       req.IdempotencyKey,
		StartedAt:            time.Now(),
	}
	if err := s.jobRepo.Create(ctx, job); err != nil {
		return nil, fmt.Errorf("gagal memulai proses rollover: %v", err)
	}

	var promoted, retained, suspended, excluded, failed int64
	var executionErrors []string
	now := time.Now()

	for _, dec := range req.Decisions {
		stOID, err := bson.ObjectIDFromHex(dec.StudentID)
		if err != nil {
			failed++
			executionErrors = append(executionErrors, fmt.Sprintf("ID siswa tidak valid: %s", dec.StudentID))
			continue
		}

		student, err := s.studentRepo.FindByID(ctx, stOID)
		if err != nil || student == nil {
			failed++
			executionErrors = append(executionErrors, fmt.Sprintf("Siswa %s tidak ditemukan", dec.StudentID))
			continue
		}

		sourceClassID := student.ClassID

		switch dec.Action {
		case "promoted", "retained":
			targetClassOID, err := bson.ObjectIDFromHex(dec.TargetClassID)
			if err != nil {
				failed++
				executionErrors = append(executionErrors, fmt.Sprintf("Kelas tujuan tidak valid untuk siswa %s", student.NIS))
				continue
			}

			// Validate class belongs to school and target academic year
			targetClass, err := s.classRepo.FindByID(ctx, targetClassOID)
			if err != nil || targetClass == nil || targetClass.SchoolID != schoolID || targetClass.AcademicYearID != targetOID {
				failed++
				executionErrors = append(executionErrors, fmt.Sprintf("Kelas tujuan tidak valid atau bukan milik tahun ajaran tujuan untuk siswa %s", student.NIS))
				continue
			}

			// 1. Update Student active class
			if err := s.studentRepo.Update(ctx, stOID, bson.M{
				"classId":   targetClassOID,
				"schoolId":  schoolID,
				"status":    "active",
				"updatedAt": now,
			}); err != nil {
				failed++
				executionErrors = append(executionErrors, fmt.Sprintf("Gagal memperbarui kelas siswa %s: %v", student.NIS, err))
				continue
			}

			// 2. Create StudentEnrollment
			enrollStatus := model.EnrollmentStatusPromoted
			decStatus := model.DecisionStatusPromoted
			if dec.Action == "retained" {
				enrollStatus = model.EnrollmentStatusRetained
				decStatus = model.DecisionStatusRetained
				retained++
			} else {
				promoted++
			}

			_ = s.enrollmentRepo.Create(ctx, &model.StudentEnrollment{
				SchoolID:       schoolID,
				StudentID:      stOID,
				ClassID:        targetClassOID,
				AcademicYearID: targetOID,
				Status:         enrollStatus,
			})

			// 3. Record AcademicDecision
			_ = s.decisionRepo.Create(ctx, &model.AcademicDecision{
				SchoolID:             schoolID,
				StudentID:            stOID,
				SourceAcademicYearID: sourceOID,
				TargetAcademicYearID: &targetOID,
				SourceClassID:        &sourceClassID,
				TargetClassID:        &targetClassOID,
				ProcessType:          model.ProcessTypePromotion,
				EligibilityStatus:    model.EligibilityEligible,
				DecisionStatus:       decStatus,
				ReasonCode:           dec.ReasonCode,
				Reason:               dec.Reason,
				ReviewStatus:         model.ReviewStatusResolved,
				DecidedBy:            userID,
				DecidedAt:            now,
				JobID:                &job.ID,
			})

		case "suspended":
			suspended++
			reason := dec.Reason
			if reason == "" {
				reason = "Kenaikan kelas ditangguhkan"
			}
			_ = s.decisionRepo.Create(ctx, &model.AcademicDecision{
				SchoolID:             schoolID,
				StudentID:            stOID,
				SourceAcademicYearID: sourceOID,
				TargetAcademicYearID: &targetOID,
				SourceClassID:        &sourceClassID,
				ProcessType:          model.ProcessTypePromotion,
				EligibilityStatus:    model.EligibilityPendingReview,
				DecisionStatus:       model.DecisionStatusSuspended,
				ReasonCode:           dec.ReasonCode,
				Reason:               reason,
				ReviewStatus:         model.ReviewStatusPending,
				DecidedBy:            userID,
				DecidedAt:            now,
				JobID:                &job.ID,
			})

		case "excluded":
			excluded++
			_ = s.decisionRepo.Create(ctx, &model.AcademicDecision{
				SchoolID:             schoolID,
				StudentID:            stOID,
				SourceAcademicYearID: sourceOID,
				TargetAcademicYearID: &targetOID,
				SourceClassID:        &sourceClassID,
				ProcessType:          model.ProcessTypePromotion,
				EligibilityStatus:    model.EligibilityEligible,
				DecisionStatus:       model.DecisionStatusExcluded,
				ReasonCode:           dec.ReasonCode,
				Reason:               dec.Reason,
				ReviewStatus:         model.ReviewStatusResolved,
				DecidedBy:            userID,
				DecidedAt:            now,
				JobID:                &job.ID,
			})
		}
	}

	jobFinalStatus := model.JobStatusCompleted
	if failed > 0 {
		if promoted+retained+suspended+excluded == 0 {
			jobFinalStatus = model.JobStatusFailed
		} else {
			jobFinalStatus = model.JobStatusCompletedWithErrors
		}
	}

	summary := model.RolloverJobSummary{
		TotalStudents: int64(len(req.Decisions)),
		Promoted:      promoted,
		Retained:      retained,
		Suspended:     suspended,
		Excluded:      excluded,
		Failed:        failed,
	}

	errMsg := ""
	if len(executionErrors) > 0 {
		errMsg = strings.Join(executionErrors, "; ")
	}

	compAt := time.Now()
	_ = s.jobRepo.Update(ctx, schoolID, job.ID, bson.M{
		"status":       jobFinalStatus,
		"summary":      summary,
		"errorMessage": errMsg,
		"completedAt":  compAt,
	})

	if s.auditRepo != nil {
		_ = s.auditRepo.Log(ctx, schoolID, userID, "promotion.execute", "rollover_job", job.ID.Hex(), bson.M{
			"sourceYearId": sourceOID.Hex(),
			"targetYearId": targetOID.Hex(),
			"promoted":     promoted,
			"retained":     retained,
			"suspended":    suspended,
			"failed":       failed,
		})
	}

	return &dto.ProcessExecutionResponse{
		JobID:  job.ID.Hex(),
		Status: string(jobFinalStatus),
		Summary: dto.RolloverJobSummaryDTO{
			TotalStudents: summary.TotalStudents,
			Promoted:      summary.Promoted,
			Retained:      summary.Retained,
			Suspended:     summary.Suspended,
			Excluded:      summary.Excluded,
			Failed:        summary.Failed,
		},
		ExecutedAt: compAt,
		Errors:     executionErrors,
	}, nil
}

// GetPromotionHistory returns execution history of promotions for a year
func (s *AcademicYearService) GetPromotionHistory(ctx context.Context, schoolID, sourceYearID bson.ObjectID) ([]dto.RolloverJobResponse, error) {
	jobs, err := s.jobRepo.FindByYearAndType(ctx, schoolID, sourceYearID, model.ProcessTypePromotion)
	if err != nil {
		return nil, err
	}
	resp := make([]dto.RolloverJobResponse, 0, len(jobs))
	for _, j := range jobs {
		resp = append(resp, s.formatJobResponse(ctx, j))
	}
	return resp, nil
}

// ============================================================================
// GRADUATION
// ============================================================================

// GetGraduationCandidates loads potential graduates (final grade students)
func (s *AcademicYearService) GetGraduationCandidates(ctx context.Context, schoolID, academicYearID bson.ObjectID) ([]dto.GraduationCandidateItem, error) {
	_, err := s.ayRepo.FindByID(ctx, schoolID, academicYearID)
	if err != nil {
		return nil, errors.New("tahun ajaran tidak ditemukan")
	}

	// 1. Find classes in this year
	classes, _, err := s.classRepo.List(ctx, bson.M{
		"schoolId":       schoolID,
		"academicYearId": academicYearID,
		"status":         bson.M{"$ne": "archived"},
	}, 1, 1000)
	if err != nil {
		return nil, err
	}

	classMap := make(map[bson.ObjectID]*model.Class)
	finalClassIDs := make([]bson.ObjectID, 0)
	for _, c := range classes {
		classMap[c.ID] = c
		// Identify final grade classes (e.g. Tingkat 6, 9, 12 or containing "6", "9", "12", "XII", "IX", "VI")
		if isFinalGradeClass(c) {
			finalClassIDs = append(finalClassIDs, c.ID)
		}
	}

	// If no specific final classes found, allow candidates from all classes in this year to be selectable
	if len(finalClassIDs) == 0 {
		for _, c := range classes {
			finalClassIDs = append(finalClassIDs, c.ID)
		}
	}

	if len(finalClassIDs) == 0 {
		return []dto.GraduationCandidateItem{}, nil
	}

	// 2. Active students
	studentFilter := bson.M{
		"classId": bson.M{"$in": finalClassIDs},
		"status":  "active",
	}
	students, _, err := s.studentRepo.List(ctx, studentFilter, 1, 5000)
	if err != nil {
		return nil, err
	}

	// 3. User map
	userIDs := make([]bson.ObjectID, 0)
	for _, st := range students {
		if !st.UserID.IsZero() {
			userIDs = append(userIDs, st.UserID)
		}
	}
	userMap := make(map[bson.ObjectID]*model.User)
	if len(userIDs) > 0 {
		cur, err := s.db.Collection("users").Find(ctx, bson.M{"_id": bson.M{"$in": userIDs}})
		if err == nil {
			var users []model.User
			if err := cur.All(ctx, &users); err == nil {
				for i := range users {
					userMap[users[i].ID] = &users[i]
				}
			}
		}
	}

	candidates := make([]dto.GraduationCandidateItem, 0, len(students))
	for _, st := range students {
		cls := classMap[st.ClassID]
		clsName := "-"
		eduLevel := "-"
		if cls != nil {
			clsName = cls.Name
			eduLevel = cls.EducationLevel
		}

		studentName := fmt.Sprintf("Siswa %s", st.NIS)
		if u, exists := userMap[st.UserID]; exists && u.Name != "" {
			studentName = u.Name
		}

		avgScore, missingGrades := s.evaluateStudentGrades(ctx, st.ID, st.ClassID, academicYearID)

		eligibility := "eligible"
		recommended := "graduated"
		reasonCode := ""
		reason := ""

		if missingGrades < 0 {
			eligibility = "pending_review"
			recommended = "suspended"
			reasonCode = "no_grades"
			reason = "Belum ada nilai"
		} else if missingGrades > 0 {
			eligibility = "pending_review"
			recommended = "suspended"
			reasonCode = "graduation_requirements_incomplete"
			reason = fmt.Sprintf("Persyaratan kelulusan belum lengkap (%d nilai belum masuk)", missingGrades)
		} else if avgScore < 65.0 && avgScore > 0 {
			eligibility = "pending_review"
			recommended = "suspended"
			reasonCode = "academic_below_graduation_threshold"
			reason = fmt.Sprintf("Rata-rata nilai (%.1f) di bawah batas kelulusan standar", avgScore)
		}

		candidates = append(candidates, dto.GraduationCandidateItem{
			StudentID:          st.ID.Hex(),
			NIS:                st.NIS,
			Name:               studentName,
			ClassID:            st.ClassID.Hex(),
			ClassName:          clsName,
			EducationLevel:     eduLevel,
			EnrollmentYear:     st.EnrollmentYear,
			EligibilityStatus:  eligibility,
			RecommendedAction:  recommended,
			ReasonCode:         reasonCode,
			Reason:             reason,
			AverageScore:       avgScore,
			MissingGradesCount: missingGrades,
		})
	}

	return candidates, nil
}

// PreviewGraduation previews graduation decisions
func (s *AcademicYearService) PreviewGraduation(ctx context.Context, schoolID bson.ObjectID, req dto.GraduationPreviewRequest) (*dto.GraduationPreviewResponse, error) {
	ayOID, err := bson.ObjectIDFromHex(req.AcademicYearID)
	if err != nil {
		return nil, errors.New("id tahun ajaran tidak valid")
	}

	candidates, err := s.GetGraduationCandidates(ctx, schoolID, ayOID)
	if err != nil {
		return nil, err
	}

	decMap := make(map[string]dto.GraduationDecisionInput)
	for _, d := range req.Decisions {
		decMap[d.StudentID] = d
	}

	var eligibleCount, pendingCount, graduatedCount, suspendedCount, excludedCount int64
	var previewErrors []string

	for i := range candidates {
		c := &candidates[i]
		if d, ok := decMap[c.StudentID]; ok {
			c.RecommendedAction = d.Action
			if d.Reason != "" {
				c.Reason = d.Reason
			}
			if d.ReasonCode != "" {
				c.ReasonCode = d.ReasonCode
			}
		}

		switch c.RecommendedAction {
		case "graduated":
			graduatedCount++
		case "suspended":
			suspendedCount++
			if c.Reason == "" {
				c.Reason = "Kelulusan ditangguhkan untuk peninjauan persyaratan"
			}
		case "excluded":
			excludedCount++
		}

		if c.EligibilityStatus == "eligible" {
			eligibleCount++
		} else {
			pendingCount++
		}
	}

	return &dto.GraduationPreviewResponse{
		TotalStudents: int64(len(candidates)),
		Eligible:      eligibleCount,
		PendingReview: pendingCount,
		Graduated:     graduatedCount,
		Suspended:     suspendedCount,
		Excluded:      excludedCount,
		Errors:        previewErrors,
		Candidates:    candidates,
	}, nil
}

// ExecuteGraduation mass-graduates students safely
func (s *AcademicYearService) ExecuteGraduation(ctx context.Context, schoolID, userID bson.ObjectID, req dto.GraduationExecutionRequest) (*dto.ProcessExecutionResponse, error) {
	if req.IdempotencyKey == "" {
		return nil, errors.New("idempotencyKey wajib disertakan")
	}

	// 1. Idempotency Check
	existingJob, err := s.jobRepo.FindByIdempotencyKey(ctx, schoolID, req.IdempotencyKey)
	if err == nil && existingJob != nil {
		return &dto.ProcessExecutionResponse{
			JobID: existingJob.ID.Hex(),
			Status: string(existingJob.Status),
			Summary: dto.RolloverJobSummaryDTO{
				TotalStudents: existingJob.Summary.TotalStudents,
				Graduated:     existingJob.Summary.Graduated,
				Suspended:     existingJob.Summary.Suspended,
				Excluded:      existingJob.Summary.Excluded,
				Failed:        existingJob.Summary.Failed,
			},
			ExecutedAt: existingJob.CreatedAt,
		}, nil
	}

	ayOID, err := bson.ObjectIDFromHex(req.AcademicYearID)
	if err != nil {
		return nil, errors.New("id tahun ajaran tidak valid")
	}

	ay, err := s.ayRepo.FindByID(ctx, schoolID, ayOID)
	if err != nil {
		return nil, errors.New("tahun ajaran tidak ditemukan")
	}

	job := &model.RolloverJob{
		SchoolID:             schoolID,
		SourceAcademicYearID: ayOID,
		ProcessType:          model.ProcessTypeGraduation,
		Status:               model.JobStatusProcessing,
		ExecutedBy:           userID,
		IdempotencyKey:       req.IdempotencyKey,
		StartedAt:            time.Now(),
	}
	if err := s.jobRepo.Create(ctx, job); err != nil {
		return nil, fmt.Errorf("gagal memulai proses kelulusan: %v", err)
	}

	var graduated, suspended, excluded, failed int64
	var executionErrors []string
	now := time.Now()

	// Extract graduation year integer from name (e.g. "2026/2027" -> 2027)
	gradYear := now.Year()
	parts := strings.Split(ay.Name, "/")
	if len(parts) == 2 {
		if y, err := strconv.Atoi(strings.TrimSpace(parts[1])); err == nil {
			gradYear = y
		}
	}

	for _, dec := range req.Decisions {
		stOID, err := bson.ObjectIDFromHex(dec.StudentID)
		if err != nil {
			failed++
			executionErrors = append(executionErrors, fmt.Sprintf("ID siswa tidak valid: %s", dec.StudentID))
			continue
		}

		student, err := s.studentRepo.FindByID(ctx, stOID)
		if err != nil || student == nil {
			failed++
			executionErrors = append(executionErrors, fmt.Sprintf("Siswa %s tidak ditemukan", dec.StudentID))
			continue
		}

		if student.Status == "graduated" {
			// Already graduated, skip safely
			graduated++
			continue
		}

		sourceClassID := student.ClassID

		switch dec.Action {
		case "graduated":
			gradDate := now
			if dec.GraduationDate != nil && !dec.GraduationDate.IsZero() {
				gradDate = *dec.GraduationDate
			}

			// Update student status to graduated, preserve all other fields
			if err := s.studentRepo.Update(ctx, stOID, bson.M{
				"status":         "graduated",
				"graduationYear": gradYear,
				"graduatedAt":    gradDate,
				"schoolId":       schoolID,
				"updatedAt":      now,
			}); err != nil {
				failed++
				executionErrors = append(executionErrors, fmt.Sprintf("Gagal mengubah status kelulusan siswa %s: %v", student.NIS, err))
				continue
			}

			graduated++

			// Record in StudentEnrollment as completed/graduated
			_ = s.enrollmentRepo.Create(ctx, &model.StudentEnrollment{
				SchoolID:       schoolID,
				StudentID:      stOID,
				ClassID:        sourceClassID,
				AcademicYearID: ayOID,
				Status:         model.EnrollmentStatusGraduated,
			})

			// Record AcademicDecision
			_ = s.decisionRepo.Create(ctx, &model.AcademicDecision{
				SchoolID:             schoolID,
				StudentID:            stOID,
				SourceAcademicYearID: ayOID,
				SourceClassID:        &sourceClassID,
				ProcessType:          model.ProcessTypeGraduation,
				EligibilityStatus:    model.EligibilityEligible,
				DecisionStatus:       model.DecisionStatusGraduated,
				ReasonCode:           dec.ReasonCode,
				Reason:               dec.Reason,
				ReviewStatus:         model.ReviewStatusResolved,
				DecidedBy:            userID,
				DecidedAt:            now,
				JobID:                &job.ID,
			})

		case "suspended":
			suspended++
			reason := dec.Reason
			if reason == "" {
				reason = "Kelulusan ditangguhkan"
			}
			_ = s.decisionRepo.Create(ctx, &model.AcademicDecision{
				SchoolID:             schoolID,
				StudentID:            stOID,
				SourceAcademicYearID: ayOID,
				SourceClassID:        &sourceClassID,
				ProcessType:          model.ProcessTypeGraduation,
				EligibilityStatus:    model.EligibilityPendingReview,
				DecisionStatus:       model.DecisionStatusSuspended,
				ReasonCode:           dec.ReasonCode,
				Reason:               reason,
				ReviewStatus:         model.ReviewStatusPending,
				DecidedBy:            userID,
				DecidedAt:            now,
				JobID:                &job.ID,
			})

		case "excluded":
			excluded++
			_ = s.decisionRepo.Create(ctx, &model.AcademicDecision{
				SchoolID:             schoolID,
				StudentID:            stOID,
				SourceAcademicYearID: ayOID,
				SourceClassID:        &sourceClassID,
				ProcessType:          model.ProcessTypeGraduation,
				EligibilityStatus:    model.EligibilityEligible,
				DecisionStatus:       model.DecisionStatusExcluded,
				ReasonCode:           dec.ReasonCode,
				Reason:               dec.Reason,
				ReviewStatus:         model.ReviewStatusResolved,
				DecidedBy:            userID,
				DecidedAt:            now,
				JobID:                &job.ID,
			})
		}
	}

	jobFinalStatus := model.JobStatusCompleted
	if failed > 0 {
		if graduated+suspended+excluded == 0 {
			jobFinalStatus = model.JobStatusFailed
		} else {
			jobFinalStatus = model.JobStatusCompletedWithErrors
		}
	}

	summary := model.RolloverJobSummary{
		TotalStudents: int64(len(req.Decisions)),
		Graduated:     graduated,
		Suspended:     suspended,
		Excluded:      excluded,
		Failed:        failed,
	}

	errMsg := ""
	if len(executionErrors) > 0 {
		errMsg = strings.Join(executionErrors, "; ")
	}

	compAt := time.Now()
	_ = s.jobRepo.Update(ctx, schoolID, job.ID, bson.M{
		"status":       jobFinalStatus,
		"summary":      summary,
		"errorMessage": errMsg,
		"completedAt":  compAt,
	})

	if s.auditRepo != nil {
		_ = s.auditRepo.Log(ctx, schoolID, userID, "graduation.execute", "rollover_job", job.ID.Hex(), bson.M{
			"academicYearId": ayOID.Hex(),
			"graduated":      graduated,
			"suspended":      suspended,
			"failed":         failed,
		})
	}

	return &dto.ProcessExecutionResponse{
		JobID:  job.ID.Hex(),
		Status: string(jobFinalStatus),
		Summary: dto.RolloverJobSummaryDTO{
			TotalStudents: summary.TotalStudents,
			Graduated:     summary.Graduated,
			Suspended:     summary.Suspended,
			Excluded:      summary.Excluded,
			Failed:        summary.Failed,
		},
		ExecutedAt: compAt,
		Errors:     executionErrors,
	}, nil
}

// GetGraduationHistory returns historical graduation runs
func (s *AcademicYearService) GetGraduationHistory(ctx context.Context, schoolID, academicYearID bson.ObjectID) ([]dto.RolloverJobResponse, error) {
	jobs, err := s.jobRepo.FindByYearAndType(ctx, schoolID, academicYearID, model.ProcessTypeGraduation)
	if err != nil {
		return nil, err
	}
	resp := make([]dto.RolloverJobResponse, 0, len(jobs))
	for _, j := range jobs {
		resp = append(resp, s.formatJobResponse(ctx, j))
	}
	return resp, nil
}

// ============================================================================
// SUSPENSION & DECISION MANAGEMENT
// ============================================================================

// GetSuspendedDecisions retrieves all unresolved suspended decisions for a given academic year
func (s *AcademicYearService) GetSuspendedDecisions(ctx context.Context, schoolID, academicYearID bson.ObjectID) ([]dto.SuspendedDecisionItem, error) {
	decisions, err := s.decisionRepo.FindSuspendedByYear(ctx, schoolID, academicYearID)
	if err != nil {
		return nil, err
	}

	resp := make([]dto.SuspendedDecisionItem, 0, len(decisions))
	for _, d := range decisions {
		stName := "-"
		stNIS := "-"
		clsName := "-"

		student, err := s.studentRepo.FindByID(ctx, d.StudentID)
		if err == nil && student != nil {
			stNIS = student.NIS
			if !student.UserID.IsZero() {
				u, err := s.userRepo.FindByID(ctx, student.UserID)
				if err == nil && u != nil {
					stName = u.Name
				}
			}
			if stName == "-" {
				stName = fmt.Sprintf("Siswa %s", student.NIS)
			}
		}

		if d.SourceClassID != nil && !d.SourceClassID.IsZero() {
			c, err := s.classRepo.FindByID(ctx, *d.SourceClassID)
			if err == nil && c != nil {
				clsName = c.Name
			}
		}

		reviewerName := ""
		if d.ReviewedBy != nil && !d.ReviewedBy.IsZero() {
			u, err := s.userRepo.FindByID(ctx, *d.ReviewedBy)
			if err == nil && u != nil {
				reviewerName = u.Name
			}
		}

		targetYearHex := ""
		if d.TargetAcademicYearID != nil {
			targetYearHex = d.TargetAcademicYearID.Hex()
		}

		resp = append(resp, dto.SuspendedDecisionItem{
			ID:                   d.ID.Hex(),
			StudentID:            d.StudentID.Hex(),
			StudentName:          stName,
			NIS:                  stNIS,
			SourceClassName:      clsName,
			SourceAcademicYearID: d.SourceAcademicYearID.Hex(),
			TargetAcademicYearID: targetYearHex,
			ProcessType:          string(d.ProcessType),
			EligibilityStatus:    string(d.EligibilityStatus),
			DecisionStatus:       string(d.DecisionStatus),
			ReasonCode:           d.ReasonCode,
			Reason:               d.Reason,
			ReviewStatus:         string(d.ReviewStatus),
			ReviewedByName:       reviewerName,
			ReviewedAt:           d.ReviewedAt,
			DecidedAt:            d.DecidedAt,
			CreatedAt:            d.CreatedAt,
		})
	}

	return resp, nil
}

// GetDecisionDetail fetches a specific decision
func (s *AcademicYearService) GetDecisionDetail(ctx context.Context, schoolID, decisionID bson.ObjectID) (*dto.SuspendedDecisionItem, error) {
	d, err := s.decisionRepo.FindByID(ctx, schoolID, decisionID)
	if err != nil {
		return nil, errors.New("keputusan akademik tidak ditemukan")
	}

	stName := "-"
	stNIS := "-"
	clsName := "-"

	student, err := s.studentRepo.FindByID(ctx, d.StudentID)
	if err == nil && student != nil {
		stNIS = student.NIS
		if !student.UserID.IsZero() {
			u, err := s.userRepo.FindByID(ctx, student.UserID)
			if err == nil && u != nil {
				stName = u.Name
			}
		}
		if stName == "-" {
			stName = fmt.Sprintf("Siswa %s", student.NIS)
		}
	}

	if d.SourceClassID != nil && !d.SourceClassID.IsZero() {
		c, err := s.classRepo.FindByID(ctx, *d.SourceClassID)
		if err == nil && c != nil {
			clsName = c.Name
		}
	}

	reviewerName := ""
	if d.ReviewedBy != nil && !d.ReviewedBy.IsZero() {
		u, err := s.userRepo.FindByID(ctx, *d.ReviewedBy)
		if err == nil && u != nil {
			reviewerName = u.Name
		}
	}

	targetYearHex := ""
	if d.TargetAcademicYearID != nil {
		targetYearHex = d.TargetAcademicYearID.Hex()
	}

	return &dto.SuspendedDecisionItem{
		ID:                   d.ID.Hex(),
		StudentID:            d.StudentID.Hex(),
		StudentName:          stName,
		NIS:                  stNIS,
		SourceClassName:      clsName,
		SourceAcademicYearID: d.SourceAcademicYearID.Hex(),
		TargetAcademicYearID: targetYearHex,
		ProcessType:          string(d.ProcessType),
		EligibilityStatus:    string(d.EligibilityStatus),
		DecisionStatus:       string(d.DecisionStatus),
		ReasonCode:           d.ReasonCode,
		Reason:               d.Reason,
		ReviewStatus:         string(d.ReviewStatus),
		ReviewedByName:       reviewerName,
		ReviewedAt:           d.ReviewedAt,
		DecidedAt:            d.DecidedAt,
		CreatedAt:            d.CreatedAt,
	}, nil
}

// ReviewDecision records review progress and notes without finalizing the decision
func (s *AcademicYearService) ReviewDecision(ctx context.Context, schoolID, userID, decisionID bson.ObjectID, req dto.ReviewDecisionRequest) error {
	d, err := s.decisionRepo.FindByID(ctx, schoolID, decisionID)
	if err != nil {
		return errors.New("keputusan akademik tidak ditemukan")
	}

	now := time.Now()
	update := bson.M{
		"reviewStatus": model.ReviewStatusReviewed,
		"reviewedBy":   userID,
		"reviewedAt":   now,
		"reason":       fmt.Sprintf("%s (Catatan Peninjauan: %s)", d.Reason, req.ReviewNotes),
	}

	if err := s.decisionRepo.Update(ctx, schoolID, decisionID, update); err != nil {
		return err
	}

	if s.auditRepo != nil {
		_ = s.auditRepo.Log(ctx, schoolID, userID, "decision.review", "academic_decision", decisionID.Hex(), bson.M{
			"notes": req.ReviewNotes,
		})
	}

	return nil
}

// ResolveDecision resolves a suspended decision into promoted, retained, graduated, or excluded
func (s *AcademicYearService) ResolveDecision(ctx context.Context, schoolID, userID, decisionID bson.ObjectID, req dto.ResolveDecisionRequest) error {
	d, err := s.decisionRepo.FindByID(ctx, schoolID, decisionID)
	if err != nil {
		return errors.New("keputusan akademik tidak ditemukan")
	}

	student, err := s.studentRepo.FindByID(ctx, d.StudentID)
	if err != nil || student == nil {
		return errors.New("siswa tidak ditemukan")
	}

	now := time.Now()
	switch req.Action {
	case "promoted", "retained":
		if req.TargetClassID == "" {
			return errors.New("targetClassId wajib diisi untuk penetapan naik atau tinggal kelas")
		}
		targetClassOID, err := bson.ObjectIDFromHex(req.TargetClassID)
		if err != nil {
			return errors.New("targetClassId tidak valid")
		}
		targetClass, err := s.classRepo.FindByID(ctx, targetClassOID)
		if err != nil || targetClass == nil || targetClass.SchoolID != schoolID {
			return errors.New("kelas tujuan tidak valid dalam sekolah ini")
		}

		targetYearID := targetClass.AcademicYearID
		if d.TargetAcademicYearID != nil && !d.TargetAcademicYearID.IsZero() {
			targetYearID = *d.TargetAcademicYearID
		}

		// Update student
		if err := s.studentRepo.Update(ctx, student.ID, bson.M{
			"classId":   targetClassOID,
			"schoolId":  schoolID,
			"status":    "active",
			"updatedAt": now,
		}); err != nil {
			return err
		}

		enrollStatus := model.EnrollmentStatusPromoted
		if req.Action == "retained" {
			enrollStatus = model.EnrollmentStatusRetained
		}
		_ = s.enrollmentRepo.Create(ctx, &model.StudentEnrollment{
			SchoolID:       schoolID,
			StudentID:      student.ID,
			ClassID:        targetClassOID,
			AcademicYearID: targetYearID,
			Status:         enrollStatus,
		})

		update := bson.M{
			"decisionStatus": model.DecisionStatus(req.Action),
			"reviewStatus":   model.ReviewStatusResolved,
			"targetClassId":  targetClassOID,
			"reviewedBy":     userID,
			"reviewedAt":     now,
			"reason":         fmt.Sprintf("Diselesaikan: %s", req.Notes),
		}
		_ = s.decisionRepo.Update(ctx, schoolID, decisionID, update)

	case "graduated":
		gradDate := now
		if req.GraduationDate != nil && !req.GraduationDate.IsZero() {
			gradDate = *req.GraduationDate
		}
		if err := s.studentRepo.Update(ctx, student.ID, bson.M{
			"status":         "graduated",
			"graduationYear": gradDate.Year(),
			"graduatedAt":    gradDate,
			"schoolId":       schoolID,
			"updatedAt":      now,
		}); err != nil {
			return err
		}

		_ = s.enrollmentRepo.Create(ctx, &model.StudentEnrollment{
			SchoolID:       schoolID,
			StudentID:      student.ID,
			ClassID:        student.ClassID,
			AcademicYearID: d.SourceAcademicYearID,
			Status:         model.EnrollmentStatusGraduated,
		})

		update := bson.M{
			"decisionStatus": model.DecisionStatusGraduated,
			"reviewStatus":   model.ReviewStatusResolved,
			"reviewedBy":     userID,
			"reviewedAt":     now,
			"reason":         fmt.Sprintf("Kelulusan disahkan: %s", req.Notes),
		}
		_ = s.decisionRepo.Update(ctx, schoolID, decisionID, update)

	case "excluded":
		update := bson.M{
			"decisionStatus": model.DecisionStatusExcluded,
			"reviewStatus":   model.ReviewStatusResolved,
			"reviewedBy":     userID,
			"reviewedAt":     now,
			"reason":         fmt.Sprintf("Dikecualikan: %s", req.Notes),
		}
		_ = s.decisionRepo.Update(ctx, schoolID, decisionID, update)
	}

	if s.auditRepo != nil {
		_ = s.auditRepo.Log(ctx, schoolID, userID, "decision.resolve", "academic_decision", decisionID.Hex(), bson.M{
			"action":    req.Action,
			"studentId": student.ID.Hex(),
			"notes":     req.Notes,
		})
	}

	return nil
}

// ============================================================================
// HELPER METHODS
// ============================================================================

func (s *AcademicYearService) evaluateStudentGrades(ctx context.Context, studentID, classID, academicYearID bson.ObjectID) (float64, int) {
	if s.db == nil {
		return 0.0, -1
	}

	filter := bson.M{
		"studentId":      studentID,
		"academicYearId": academicYearID,
	}
	cur, err := s.db.Collection("grades").Find(ctx, filter)
	if err != nil {
		return 0.0, 0
	}
	defer cur.Close(ctx)

	var grades []model.Grade
	if err := cur.All(ctx, &grades); err != nil || len(grades) == 0 {
		// Belum ada nilai
		return 0.0, -1
	}

	var sum float64
	for _, g := range grades {
		sum += g.Score
	}
	avg := sum / float64(len(grades))

	return avg, 0
}

func (s *AcademicYearService) formatJobResponse(ctx context.Context, j *model.RolloverJob) dto.RolloverJobResponse {
	srcName := "-"
	targetName := "-"
	targetYearHex := ""
	execName := "Administrator"

	if s.ayRepo != nil {
		if y, err := s.ayRepo.FindByID(ctx, j.SchoolID, j.SourceAcademicYearID); err == nil && y != nil {
			srcName = y.Name
		}
		if j.TargetAcademicYearID != nil && !j.TargetAcademicYearID.IsZero() {
			targetYearHex = j.TargetAcademicYearID.Hex()
			if y, err := s.ayRepo.FindByID(ctx, j.SchoolID, *j.TargetAcademicYearID); err == nil && y != nil {
				targetName = y.Name
			}
		}
	}

	if s.userRepo != nil && !j.ExecutedBy.IsZero() {
		if u, err := s.userRepo.FindByID(ctx, j.ExecutedBy); err == nil && u != nil {
			execName = u.Name
		}
	}

	return dto.RolloverJobResponse{
		ID:                     j.ID.Hex(),
		ProcessType:            string(j.ProcessType),
		SourceAcademicYearID:   j.SourceAcademicYearID.Hex(),
		SourceAcademicYearName: srcName,
		TargetAcademicYearID:   targetYearHex,
		TargetAcademicYearName: targetName,
		Status:                 string(j.Status),
		Summary: dto.RolloverJobSummaryDTO{
			TotalStudents: j.Summary.TotalStudents,
			Promoted:      j.Summary.Promoted,
			Retained:      j.Summary.Retained,
			Graduated:     j.Summary.Graduated,
			Suspended:     j.Summary.Suspended,
			Excluded:      j.Summary.Excluded,
			Failed:        j.Summary.Failed,
		},
		ExecutedByName: execName,
		ErrorMessage:   j.ErrorMessage,
		StartedAt:      j.StartedAt,
		CompletedAt:    j.CompletedAt,
	}
}

var (
	romanFinalGradeRegex  = regexp.MustCompile(`(?i)(?:^|[\s\-_/])(XII|IX|VI)(?:$|[\s\-_/])`)
	arabicFinalGradeRegex = regexp.MustCompile(`(?i)(?:^|[\s\-_/])(12|9|6)(?:$|[\s\-_/]|[A-Z])`)
)

func isFinalGradeClass(c *model.Class) bool {
	if c == nil {
		return false
	}
	name := strings.TrimSpace(c.Name)
	code := strings.TrimSpace(c.Code)

	if romanFinalGradeRegex.MatchString(name) || romanFinalGradeRegex.MatchString(code) {
		return true
	}
	if arabicFinalGradeRegex.MatchString(name) || arabicFinalGradeRegex.MatchString(code) {
		return true
	}
	return false
}

