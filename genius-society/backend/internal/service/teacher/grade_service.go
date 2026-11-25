package teachersvc

import (
	"context"
	"math"
	"time"

	teacherdto "github.com/geniussociety/backend/internal/dto/teacher"
	"github.com/geniussociety/backend/internal/handler"
	"github.com/geniussociety/backend/internal/middleware"
	"github.com/geniussociety/backend/internal/model"
	teacherrepo "github.com/geniussociety/backend/internal/repository/teacher"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// GradeService defines operations for managing and retrieving student grades.
type GradeService interface {
	GetGrades(ctx context.Context, scope *middleware.TeacherScope, classID, subjectID bson.ObjectID) ([]teacherdto.StudentGradeRowDTO, error)
	UpdateGrade(ctx context.Context, scope *middleware.TeacherScope, payload teacherdto.UpdateStudentGradePayload) (*teacherdto.GradeItemDTO, error)

	InputSTSExams(ctx context.Context, scope *middleware.TeacherScope, payload *teacherdto.BulkExamsPayload) error
	SetManualGrade(ctx context.Context, scope *middleware.TeacherScope, itemType string, itemID, studentID bson.ObjectID, score float64, feedback, reason string) error
	DeleteManualGrade(ctx context.Context, scope *middleware.TeacherScope, itemType string, itemID, studentID bson.ObjectID) error
	GetGradeRecap(ctx context.Context, scope *middleware.TeacherScope, classID, subjectID bson.ObjectID, semester int) (*teacherdto.GradeRecapResponseDTO, error)
}

type gradeServiceImpl struct {
	repo teacherrepo.GradeRepository
}

// NewGradeService creates a new GradeService instance.
func NewGradeService(repo teacherrepo.GradeRepository) GradeService {
	return &gradeServiceImpl{repo: repo}
}

func (s *gradeServiceImpl) GetGrades(ctx context.Context, scope *middleware.TeacherScope, classID, subjectID bson.ObjectID) ([]teacherdto.StudentGradeRowDTO, error) {
	// 1. Authorization: teacher must teach this class+subject, or be walas of this class
	if !scope.HasClassSubject(classID, subjectID) && !scope.IsWalasOf(classID) {
		return nil, handler.NewForbiddenError("Anda tidak memiliki akses penilaian pada kelas atau mata pelajaran ini")
	}

	// 2. Fetch students in the class
	students, err := s.repo.FindStudentsByClassID(ctx, classID)
	if err != nil {
		return nil, err
	}
	if len(students) == 0 {
		return []teacherdto.StudentGradeRowDTO{}, nil
	}

	// 3. User lookup for names
	var userIDs []bson.ObjectID
	for _, st := range students {
		if !st.UserID.IsZero() {
			userIDs = append(userIDs, st.UserID)
		}
	}
	userMap, _ := s.repo.FindUsersByIDs(ctx, userIDs)

	// 4. Fetch grades for this class & subject
	grades, err := s.repo.FindGradesByClassSubject(ctx, classID, subjectID, scope.AcademicYearID)
	if err != nil {
		return nil, err
	}

	// 5. Fetch grading configuration (default 60/20/20)
	dailyWeight := 0.6
	utsWeight := 0.2
	uasWeight := 0.2
	cfg, _ := s.repo.FindGradingConfig(ctx, scope.AcademicYearID, subjectID)
	if cfg != nil && (cfg.DailyWeight+cfg.UTSWeight+cfg.UASWeight) > 0 {
		totalW := cfg.DailyWeight + cfg.UTSWeight + cfg.UASWeight
		dailyWeight = cfg.DailyWeight / totalW
		utsWeight = cfg.UTSWeight / totalW
		uasWeight = cfg.UASWeight / totalW
	}

	// 6. Index grades by student ID and type
	type studentGradeAccum struct {
		assignmentScores []float64
		assessmentScores []float64
		projectScores    []float64
		utsScore         *float64
		uasScore         *float64
	}
	gradeMap := make(map[bson.ObjectID]*studentGradeAccum)

	for _, g := range grades {
		accum, ok := gradeMap[g.StudentID]
		if !ok {
			accum = &studentGradeAccum{}
			gradeMap[g.StudentID] = accum
		}

		switch g.Type {
		case model.GradeAssignment:
			accum.assignmentScores = append(accum.assignmentScores, g.Score)
		case model.GradeAssessment:
			accum.assessmentScores = append(accum.assessmentScores, g.Score)
		case model.GradeProject:
			accum.projectScores = append(accum.projectScores, g.Score)
		case model.GradeUTS:
			score := g.Score
			accum.utsScore = &score
		case model.GradeUAS:
			score := g.Score
			accum.uasScore = &score
		}
	}

	// 7. Build rows
	round1 := func(val float64) float64 {
		return math.Round(val*10) / 10
	}

	rows := make([]teacherdto.StudentGradeRowDTO, 0, len(students))
	for _, st := range students {
		stName := "Siswa"
		if u, ok := userMap[st.UserID]; ok && u != nil {
			stName = u.Name
		}

		accum := gradeMap[st.ID]
		var asScore, asmScore, prScore, dAvg, fScore *float64
		var uts, uas *float64

		if accum != nil {
			uts = accum.utsScore
			uas = accum.uasScore

			var dailySum float64
			var dailyCount int

			if len(accum.assignmentScores) > 0 {
				var sum float64
				for _, sc := range accum.assignmentScores {
					sum += sc
				}
				avg := round1(sum / float64(len(accum.assignmentScores)))
				asScore = &avg
				dailySum += sum
				dailyCount += len(accum.assignmentScores)
			}

			if len(accum.assessmentScores) > 0 {
				var sum float64
				for _, sc := range accum.assessmentScores {
					sum += sc
				}
				avg := round1(sum / float64(len(accum.assessmentScores)))
				asmScore = &avg
				dailySum += sum
				dailyCount += len(accum.assessmentScores)
			}

			if len(accum.projectScores) > 0 {
				var sum float64
				for _, sc := range accum.projectScores {
					sum += sc
				}
				avg := round1(sum / float64(len(accum.projectScores)))
				prScore = &avg
				dailySum += sum
				dailyCount += len(accum.projectScores)
			}

			if dailyCount > 0 {
				avg := round1(dailySum / float64(dailyCount))
				dAvg = &avg
			}

			// Calculate final score if any score exists
			var finalVal float64
			var hasAnyScore bool

			if dAvg != nil {
				finalVal += *dAvg * dailyWeight
				hasAnyScore = true
			}
			if uts != nil {
				finalVal += *uts * utsWeight
				hasAnyScore = true
			}
			if uas != nil {
				finalVal += *uas * uasWeight
				hasAnyScore = true
			}

			if hasAnyScore {
				roundedFinal := round1(finalVal)
				fScore = &roundedFinal
			}
		}

		rows = append(rows, teacherdto.StudentGradeRowDTO{
			StudentID:       st.ID.Hex(),
			StudentName:     stName,
			StudentNIS:      st.NIS,
			AssignmentScore: asScore,
			AssessmentScore: asmScore,
			ProjectScore:    prScore,
			DailyAvg:        dAvg,
			UTSScore:        uts,
			UASScore:        uas,
			FinalScore:      fScore,
		})
	}

	return rows, nil
}

func (s *gradeServiceImpl) UpdateGrade(ctx context.Context, scope *middleware.TeacherScope, payload teacherdto.UpdateStudentGradePayload) (*teacherdto.GradeItemDTO, error) {
	studentID, err := bson.ObjectIDFromHex(payload.StudentID)
	if err != nil {
		return nil, handler.NewValidationError("ID siswa tidak valid")
	}
	subjectID, err := bson.ObjectIDFromHex(payload.SubjectID)
	if err != nil {
		return nil, handler.NewValidationError("ID mata pelajaran tidak valid")
	}
	classID, err := bson.ObjectIDFromHex(payload.ClassID)
	if err != nil {
		return nil, handler.NewValidationError("ID kelas tidak valid")
	}

	// Must teach this class & subject
	if !scope.HasClassSubject(classID, subjectID) {
		return nil, handler.NewForbiddenError("Anda tidak memiliki hak akses mengubah nilai pada kelas/mata pelajaran ini")
	}

	gradeType := model.GradeType(payload.Type)
	if gradeType != model.GradeUTS && gradeType != model.GradeUAS {
		return nil, handler.NewValidationError("Tipe nilai hanya boleh 'uts' atau 'uas'")
	}

	gradeModel := &model.Grade{
		StudentID:      studentID,
		SubjectID:      subjectID,
		ClassID:        classID,
		AcademicYearID: scope.AcademicYearID,
		Semester:       scope.Semester,
		TeacherID:      scope.TeacherID,
		Type:           gradeType,
		Score:          payload.Score,
		Source:         "manual",
	}

	saved, err := s.repo.UpsertExamGrade(ctx, gradeModel)
	if err != nil {
		return nil, err
	}

	// Fetch student and user names for response
	st, _ := s.repo.FindStudentByID(ctx, studentID)
	stName := ""
	stNIS := ""
	if st != nil {
		stNIS = st.NIS
		userMap, _ := s.repo.FindUsersByIDs(ctx, []bson.ObjectID{st.UserID})
		if u, ok := userMap[st.UserID]; ok && u != nil {
			stName = u.Name
		}
	}

	cls, _ := s.repo.FindClassByID(ctx, classID)
	clsName := ""
	if cls != nil {
		clsName = cls.Name
	}

	sub, _ := s.repo.FindSubjectByID(ctx, subjectID)
	subName := ""
	if sub != nil {
		subName = sub.Name
	}

	return &teacherdto.GradeItemDTO{
		ID:             saved.ID.Hex(),
		StudentID:      saved.StudentID.Hex(),
		StudentName:    stName,
		StudentNIS:     stNIS,
		SubjectID:      saved.SubjectID.Hex(),
		SubjectName:    subName,
		ClassID:        saved.ClassID.Hex(),
		ClassName:      clsName,
		AcademicYearID: saved.AcademicYearID.Hex(),
		Semester:       saved.Semester,
		Type:           string(saved.Type),
		Score:          saved.Score,
		CreatedAt:      saved.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *gradeServiceImpl) InputSTSExams(ctx context.Context, scope *middleware.TeacherScope, payload *teacherdto.BulkExamsPayload) error {
	if scope == nil {
		return handler.ErrUnauthorized
	}

	if payload.Type != "sts" {
		return handler.NewValidationError("Tipe ujian hanya mendukung 'sts'")
	}

	if len(payload.Entries) == 0 {
		return handler.NewValidationError("Daftar nilai ujian tidak boleh kosong")
	}

	if len(payload.Entries) > 100 {
		return handler.NewValidationError("Maksimal 100 entri nilai per pengiriman")
	}

	classID, err := bson.ObjectIDFromHex(payload.ClassID)
	if err != nil {
		return handler.NewValidationError("ID kelas tidak valid")
	}

	subjectID, err := bson.ObjectIDFromHex(payload.SubjectID)
	if err != nil {
		return handler.NewValidationError("ID mata pelajaran tidak valid")
	}

	if !scope.HasClassSubject(classID, subjectID) {
		return handler.ErrNotAssigned
	}

	semester := payload.Semester
	if semester <= 0 {
		semester = scope.Semester
	}

	// Validate all students belong to the class
	students, err := s.repo.FindStudentsByClassID(ctx, classID)
	if err != nil {
		return err
	}
	studentClassMap := make(map[bson.ObjectID]bool)
	for _, st := range students {
		studentClassMap[st.ID] = true
	}

	now := time.Now().UTC()
	var grades []*model.Grade
	var audits []*model.AuditLog

	seenStudents := make(map[bson.ObjectID]bool)
	for _, entry := range payload.Entries {
		stID, err := bson.ObjectIDFromHex(entry.StudentID)
		if err != nil {
			return handler.NewValidationError("ID siswa tidak valid")
		}
		if !studentClassMap[stID] {
			return handler.NewValidationError("Siswa tidak aktif pada kelas sasaran")
		}
		if seenStudents[stID] {
			return handler.NewValidationError("Duplikasi siswa pada daftar entri nilai")
		}
		seenStudents[stID] = true

		if entry.Score < 0 || entry.Score > 100 {
			return handler.NewValidationError("Nilai harus berada di rentang 0 hingga 100")
		}

		// Max 1 STS per student-subject-semester: duplicate check
		existing, err := s.repo.FindExistingExamGrade(ctx, stID, subjectID, scope.AcademicYearID, semester, model.GradeUTS)
		if err != nil {
			return err
		}
		if existing != nil {
			return handler.NewConflictError("Nilai STS untuk siswa sudah tercatat pada semester ini")
		}

		g := &model.Grade{
			StudentID:      stID,
			SubjectID:      subjectID,
			ClassID:        classID,
			AcademicYearID: scope.AcademicYearID,
			Semester:       semester,
			TeacherID:      scope.TeacherID,
			Type:           model.GradeUTS,
			Score:          entry.Score,
			Source:         "manual",
			GradedBy:       scope.TeacherID,
		}
		grades = append(grades, g)

		audit := &model.AuditLog{
			ID:         bson.NewObjectID(),
			UserID:     scope.UserID,
			Action:     "input_sts_exam",
			EntityType: "grade",
			EntityID:   stID.Hex(),
			Details: bson.M{
				"studentId": stID.Hex(),
				"subjectId": subjectID.Hex(),
				"classId":   classID.Hex(),
				"score":     entry.Score,
				"semester":  semester,
			},
			Timestamp: now,
		}
		audits = append(audits, audit)
	}

	return s.repo.BulkUpsertSTS(ctx, scope.TeacherID, grades, audits)
}

func (s *gradeServiceImpl) SetManualGrade(ctx context.Context, scope *middleware.TeacherScope, itemType string, itemID, studentID bson.ObjectID, score float64, feedback, reason string) error {
	if scope == nil {
		return handler.ErrUnauthorized
	}

	if score < 0 || score > 100 {
		return handler.NewValidationError("Nilai harus berada di rentang 0 hingga 100")
	}

	student, err := s.repo.FindStudentByID(ctx, studentID)
	if err != nil || student == nil {
		return handler.ErrNotFound
	}

	now := time.Now().UTC()
	var subjectID bson.ObjectID
	var classIDs []bson.ObjectID
	var gradeType model.GradeType

	switch itemType {
	case "assignment":
		asg, err := s.repo.FindAssignment(ctx, scope.TeacherID, itemID)
		if err != nil || asg == nil {
			return handler.ErrNotFound
		}
		if now.Before(asg.Deadline) {
			return handler.NewConflictError("DEADLINE_NOT_PASSED: Penilaian manual hanya dapat dilakukan setelah batas waktu pengumpulan lewat")
		}
		hasSub, _ := s.repo.HasAssignmentSubmission(ctx, itemID, studentID)
		if hasSub {
			return handler.NewConflictError("SUBMISSION_EXISTS: Siswa sudah memiliki pengumpulan tugas")
		}
		subjectID = asg.SubjectID
		classIDs = asg.ClassIDs
		gradeType = model.GradeAssignment

	case "project":
		prj, err := s.repo.FindProject(ctx, scope.TeacherID, itemID)
		if err != nil || prj == nil {
			return handler.ErrNotFound
		}
		if now.Before(prj.Deadline) {
			return handler.NewConflictError("DEADLINE_NOT_PASSED: Penilaian manual hanya dapat dilakukan setelah batas waktu pengumpulan lewat")
		}
		hasSub, _ := s.repo.HasProjectSubmission(ctx, itemID, studentID)
		if hasSub {
			return handler.NewConflictError("SUBMISSION_EXISTS: Siswa sudah memiliki pengumpulan projek")
		}
		subjectID = prj.SubjectID
		classIDs = prj.ClassIDs
		gradeType = model.GradeProject

	case "assessment":
		asm, err := s.repo.FindAssessment(ctx, scope.TeacherID, itemID)
		if err != nil || asm == nil {
			return handler.ErrNotFound
		}
		if asm.EndAt != nil && now.Before(*asm.EndAt) {
			return handler.NewConflictError("DEADLINE_NOT_PASSED: Penilaian manual hanya dapat dilakukan setelah batas waktu asesmen selesai")
		}
		hasAtt, _ := s.repo.HasAssessmentAttempt(ctx, itemID, studentID)
		if hasAtt {
			return handler.NewConflictError("SUBMISSION_EXISTS: Siswa sudah memiliki pengerjaan asesmen")
		}
		subjectID = asm.SubjectID
		classIDs = asm.ClassIDs
		gradeType = model.GradeAssessment

	default:
		return handler.NewValidationError("Tipe item tidak valid")
	}

	// Verify student belongs to one of item's classes
	classMatch := false
	for _, cid := range classIDs {
		if cid == student.ClassID {
			classMatch = true
			break
		}
	}
	if !classMatch {
		return handler.ErrNotFound
	}

	grade := &model.Grade{
		StudentID:      studentID,
		SubjectID:      subjectID,
		ClassID:        student.ClassID,
		AcademicYearID: scope.AcademicYearID,
		Semester:       scope.Semester,
		TeacherID:      scope.TeacherID,
		Type:           gradeType,
		SourceID:       itemID,
		Score:          score,
		Source:         "manual",
		Reason:         reason,
		GradedBy:       scope.TeacherID,
	}

	audit := &model.AuditLog{
		ID:         bson.NewObjectID(),
		UserID:     scope.UserID,
		Action:     "manual_item_grade",
		EntityType: "grade",
		EntityID:   itemID.Hex(),
		Details: bson.M{
			"itemType":  itemType,
			"itemId":    itemID.Hex(),
			"studentId": studentID.Hex(),
			"score":     score,
			"reason":    reason,
			"feedback":  feedback,
		},
		Timestamp: now,
	}

	return s.repo.SetManualItemGrade(ctx, grade, audit)
}

func (s *gradeServiceImpl) DeleteManualGrade(ctx context.Context, scope *middleware.TeacherScope, itemType string, itemID, studentID bson.ObjectID) error {
	if scope == nil {
		return handler.ErrUnauthorized
	}

	g, err := s.repo.FindManualItemGrade(ctx, itemID, studentID)
	if err != nil || g == nil {
		return handler.ErrNotFound
	}

	if g.Source != "manual" {
		return handler.NewConflictError("Hanya nilai dengan sumber manual yang dapat dihapus")
	}

	audit := &model.AuditLog{
		ID:         bson.NewObjectID(),
		UserID:     scope.UserID,
		Action:     "delete_manual_grade",
		EntityType: "grade",
		EntityID:   g.ID.Hex(),
		Details: bson.M{
			"itemType":  itemType,
			"itemId":    itemID.Hex(),
			"studentId": studentID.Hex(),
		},
		Timestamp: time.Now().UTC(),
	}

	return s.repo.DeleteManualItemGrade(ctx, itemID, studentID, audit)
}

func (s *gradeServiceImpl) GetGradeRecap(ctx context.Context, scope *middleware.TeacherScope, classID, subjectID bson.ObjectID, semester int) (*teacherdto.GradeRecapResponseDTO, error) {
	if scope == nil {
		return nil, handler.ErrUnauthorized
	}

	if !scope.HasClassSubject(classID, subjectID) && !scope.IsWalasOf(classID) {
		return nil, handler.NewForbiddenError("Anda tidak memiliki akses rekap nilai untuk kelas dan mata pelajaran ini")
	}

	if semester <= 0 {
		semester = scope.Semester
	}

	cls, _ := s.repo.FindClassByID(ctx, classID)
	className := ""
	if cls != nil {
		className = cls.Name
	}

	sub, _ := s.repo.FindSubjectByID(ctx, subjectID)
	subjectName := ""
	if sub != nil {
		subjectName = sub.Name
	}

	// Dynamic weight resolution
	outcomeWeight := 50.0
	stsWeight := 50.0
	cfg, _ := s.repo.FindGradingConfig(ctx, scope.AcademicYearID, subjectID)
	if cfg != nil && (cfg.OutcomeWeight+cfg.StsWeight == 100) {
		outcomeWeight = cfg.OutcomeWeight
		stsWeight = cfg.StsWeight
	}

	students, err := s.repo.FindStudentsByClassID(ctx, classID)
	if err != nil {
		return nil, err
	}

	var userIDs []bson.ObjectID
	for _, st := range students {
		if !st.UserID.IsZero() {
			userIDs = append(userIDs, st.UserID)
		}
	}
	userMap, _ := s.repo.FindUsersByIDs(ctx, userIDs)

	outcomes, err := s.repo.FindLearningOutcomesForRecap(ctx, scope.AcademicYearID, subjectID, classID, semester)
	if err != nil {
		return nil, err
	}

	asgs, prjs, asms, err := s.repo.FindItemsForRecap(ctx, scope.AcademicYearID, subjectID, classID)
	if err != nil {
		return nil, err
	}

	// Map items to learning outcome
	type itemRef struct {
		id   bson.ObjectID
		kind model.GradeType
	}
	outcomeItems := make(map[bson.ObjectID][]itemRef)
	unmappedCount := 0

	for _, a := range asgs {
		if a.LearningOutcomeID != nil && !a.LearningOutcomeID.IsZero() {
			outcomeItems[*a.LearningOutcomeID] = append(outcomeItems[*a.LearningOutcomeID], itemRef{id: a.ID, kind: model.GradeAssignment})
		} else {
			unmappedCount++
		}
	}
	for _, p := range prjs {
		if p.LearningOutcomeID != nil && !p.LearningOutcomeID.IsZero() {
			outcomeItems[*p.LearningOutcomeID] = append(outcomeItems[*p.LearningOutcomeID], itemRef{id: p.ID, kind: model.GradeProject})
		} else {
			unmappedCount++
		}
	}
	for _, am := range asms {
		if am.LearningOutcomeID != nil && !am.LearningOutcomeID.IsZero() {
			outcomeItems[*am.LearningOutcomeID] = append(outcomeItems[*am.LearningOutcomeID], itemRef{id: am.ID, kind: model.GradeAssessment})
		} else {
			unmappedCount++
		}
	}

	// Fetch all grades for this class & subject & semester
	grades, err := s.repo.FindGradesByClassSubject(ctx, classID, subjectID, scope.AcademicYearID)
	if err != nil {
		return nil, err
	}

	// Index grades by studentID -> sourceID (or type)
	type studentGradeInfo struct {
		itemScores map[bson.ObjectID]float64
		manualMap  map[bson.ObjectID]bool
		stsScore   *float64
	}
	studentGrades := make(map[bson.ObjectID]*studentGradeInfo)
	for _, g := range grades {
		if g.Semester != semester {
			continue
		}
		info, ok := studentGrades[g.StudentID]
		if !ok {
			info = &studentGradeInfo{
				itemScores: make(map[bson.ObjectID]float64),
				manualMap:  make(map[bson.ObjectID]bool),
			}
			studentGrades[g.StudentID] = info
		}

		if g.Type == model.GradeUTS {
			sc := g.Score
			info.stsScore = &sc
		} else if !g.SourceID.IsZero() {
			info.itemScores[g.SourceID] = g.Score
			if g.Source == "manual" {
				info.manualMap[g.SourceID] = true
			}
		}
	}

	var emptyOutcomes []string
	studentRows := make([]teacherdto.StudentGradeRecapDTO, 0, len(students))

	for _, st := range students {
		stName := ""
		if u, ok := userMap[st.UserID]; ok && u != nil {
			stName = u.Name
		}

		info := studentGrades[st.ID]
		if info == nil {
			info = &studentGradeInfo{
				itemScores: make(map[bson.ObjectID]float64),
				manualMap:  make(map[bson.ObjectID]bool),
			}
		}

		var studentCalcOutcomes []OutcomeGradeCalculation
		var outcomeRecapDTOs []teacherdto.OutcomeRecapDTO
		var missing []string

		for _, lo := range outcomes {
			items := outcomeItems[lo.ID]
			var scores []float64
			manualCount := 0

			for _, it := range items {
				if sc, ok := info.itemScores[it.id]; ok {
					scores = append(scores, sc)
					if info.manualMap[it.id] {
						manualCount++
					}
				}
			}

			sc, prov := CalculateOutcomeGrade(scores, len(items))

			calc := OutcomeGradeCalculation{
				OutcomeID:        lo.ID.Hex(),
				Title:            lo.Title,
				ItemCount:        len(items),
				GradedCount:      len(scores),
				ManualCount:      manualCount,
				Score:            sc,
				ProvisionalScore: prov,
			}
			studentCalcOutcomes = append(studentCalcOutcomes, calc)

			if len(items) > 0 && sc == nil {
				missing = append(missing, "PP: "+lo.Title)
			}

			outcomeRecapDTOs = append(outcomeRecapDTOs, teacherdto.OutcomeRecapDTO{
				ID:               lo.ID.Hex(),
				Title:            lo.Title,
				Score:            sc,
				ProvisionalScore: prov,
				GradedCount:      len(scores),
				TotalCount:       len(items),
				ManualCount:      manualCount,
			})
		}

		rataPP, empty := CalculateRataPP(studentCalcOutcomes)
		emptyOutcomes = empty

		if info.stsScore == nil {
			missing = append(missing, "STS")
		}

		finalScore := CalculateSemester1Score(rataPP, info.stsScore, outcomeWeight, stsWeight)
		isComplete := finalScore != nil

		studentRows = append(studentRows, teacherdto.StudentGradeRecapDTO{
			StudentID:         st.ID.Hex(),
			StudentName:       stName,
			NIS:               st.NIS,
			Outcomes:          outcomeRecapDTOs,
			AvgOutcomes:       rataPP,
			STS:               info.stsScore,
			SemesterStsScore:  finalScore,
			IsComplete:        isComplete,
			MissingComponents: missing,
		})
	}

	return &teacherdto.GradeRecapResponseDTO{
		ClassID:           classID.Hex(),
		ClassName:         className,
		SubjectID:         subjectID.Hex(),
		SubjectName:       subjectName,
		Semester:          semester,
		OutcomeWeight:     outcomeWeight,
		StsWeight:         stsWeight,
		UnmappedItemCount: unmappedCount,
		EmptyOutcomes:     emptyOutcomes,
		Students:          studentRows,
	}, nil
}

