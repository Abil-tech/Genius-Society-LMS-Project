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
	students, err := s.repo.FindStudentsByClassID(ctx, scope.SchoolID, classID)
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
	grades, err := s.repo.FindGradesByClassSubject(ctx, scope.SchoolID, classID, subjectID, scope.AcademicYearID)
	if err != nil {
		return nil, err
	}

	// 5. Fetch grading configuration (default 60/20/20)
	dailyWeight := 0.6
	utsWeight := 0.2
	uasWeight := 0.2
	cfg, _ := s.repo.FindGradingConfig(ctx, scope.SchoolID, scope.AcademicYearID, subjectID)
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
		SchoolID:       scope.SchoolID,
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
