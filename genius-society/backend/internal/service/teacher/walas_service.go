package teachersvc

import (
	"context"
	"math"
	"sort"

	teacherdto "github.com/geniussociety/backend/internal/dto/teacher"
	"github.com/geniussociety/backend/internal/handler"
	"github.com/geniussociety/backend/internal/middleware"
	"github.com/geniussociety/backend/internal/model"
	teacherrepo "github.com/geniussociety/backend/internal/repository/teacher"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// WalasService defines operations exclusive to homeroom teachers (Walas).
type WalasService interface {
	GetWalasStudents(ctx context.Context, scope *middleware.TeacherScope) ([]teacherdto.WalasStudentDetailDTO, error)
	GetWalasGrades(ctx context.Context, scope *middleware.TeacherScope) ([]teacherdto.WalasGradeRekapRowDTO, error)
}

type walasServiceImpl struct {
	repo teacherrepo.WalasRepository
}

// NewWalasService constructs a new WalasService.
func NewWalasService(repo teacherrepo.WalasRepository) WalasService {
	return &walasServiceImpl{repo: repo}
}

func (s *walasServiceImpl) checkWalasScope(scope *middleware.TeacherScope) (bson.ObjectID, error) {
	if !scope.IsWalas || len(scope.WalasClassIDs) == 0 {
		return bson.NilObjectID, handler.NewForbiddenError("Akun Anda tidak tercatat sebagai Wali Kelas aktif")
	}
	return scope.WalasClassIDs[0], nil
}

func (s *walasServiceImpl) GetWalasStudents(ctx context.Context, scope *middleware.TeacherScope) ([]teacherdto.WalasStudentDetailDTO, error) {
	classID, err := s.checkWalasScope(scope)
	if err != nil {
		return nil, err
	}

	students, err := s.repo.FindStudentsByClassID(ctx, scope.SchoolID, classID)
	if err != nil {
		return nil, err
	}
	if len(students) == 0 {
		return []teacherdto.WalasStudentDetailDTO{}, nil
	}

	var userIDs []bson.ObjectID
	for _, st := range students {
		if !st.UserID.IsZero() {
			userIDs = append(userIDs, st.UserID)
		}
	}
	userMap, _ := s.repo.FindUsersByIDs(ctx, userIDs)

	// Fetch all grades to compute average score per student
	grades, _ := s.repo.FindGradesByClass(ctx, scope.SchoolID, classID, scope.AcademicYearID)
	studentScores := make(map[bson.ObjectID][]float64)
	for _, g := range grades {
		studentScores[g.StudentID] = append(studentScores[g.StudentID], g.Score)
	}

	round1 := func(v float64) float64 {
		return math.Round(v*10) / 10
	}

	results := make([]teacherdto.WalasStudentDetailDTO, 0, len(students))
	for _, st := range students {
		stName := "Siswa"
		email := ""
		if u, ok := userMap[st.UserID]; ok && u != nil {
			stName = u.Name
			email = u.Email
		}

		var avgScore *float64
		if scores, ok := studentScores[st.ID]; ok && len(scores) > 0 {
			var sum float64
			for _, sc := range scores {
				sum += sc
			}
			avg := round1(sum / float64(len(scores)))
			avgScore = &avg
		}

		results = append(results, teacherdto.WalasStudentDetailDTO{
			ID:           st.ID.Hex(),
			Name:         stName,
			NIS:          st.NIS,
			Gender:       st.Gender,
			Email:        email,
			Phone:        st.ParentPhone,
			ParentName:   st.ParentName,
			ParentPhone:  st.ParentPhone,
			AverageScore: avgScore,
		})
	}

	// Sort by averageScore descending to calculate ranks
	sort.Slice(results, func(i, j int) bool {
		s1 := float64(0)
		if results[i].AverageScore != nil {
			s1 = *results[i].AverageScore
		}
		s2 := float64(0)
		if results[j].AverageScore != nil {
			s2 = *results[j].AverageScore
		}
		return s1 > s2
	})

	for i := range results {
		results[i].Rank = i + 1
	}

	return results, nil
}

func (s *walasServiceImpl) GetWalasGrades(ctx context.Context, scope *middleware.TeacherScope) ([]teacherdto.WalasGradeRekapRowDTO, error) {
	classID, err := s.checkWalasScope(scope)
	if err != nil {
		return nil, err
	}

	students, err := s.repo.FindStudentsByClassID(ctx, scope.SchoolID, classID)
	if err != nil {
		return nil, err
	}
	if len(students) == 0 {
		return []teacherdto.WalasGradeRekapRowDTO{}, nil
	}

	var userIDs []bson.ObjectID
	for _, st := range students {
		if !st.UserID.IsZero() {
			userIDs = append(userIDs, st.UserID)
		}
	}
	userMap, _ := s.repo.FindUsersByIDs(ctx, userIDs)

	grades, _ := s.repo.FindGradesByClass(ctx, scope.SchoolID, classID, scope.AcademicYearID)

	// Collect all subject IDs from grades
	subjectIDSet := make(map[bson.ObjectID]bool)
	for _, g := range grades {
		subjectIDSet[g.SubjectID] = true
	}
	var subjectIDs []bson.ObjectID
	for sid := range subjectIDSet {
		subjectIDs = append(subjectIDs, sid)
	}
	subjectMap, _ := s.repo.FindSubjectsByIDs(ctx, subjectIDs)

	type subAccum struct {
		dailySum   float64
		dailyCount int
		uts        float64
		uas        float64
	}
	// studentID -> subjectID -> subAccum
	studentGradeMatrix := make(map[bson.ObjectID]map[bson.ObjectID]*subAccum)

	for _, g := range grades {
		sm, ok := studentGradeMatrix[g.StudentID]
		if !ok {
			sm = make(map[bson.ObjectID]*subAccum)
			studentGradeMatrix[g.StudentID] = sm
		}
		sa, ok := sm[g.SubjectID]
		if !ok {
			sa = &subAccum{}
			sm[g.SubjectID] = sa
		}

		switch g.Type {
		case model.GradeAssignment, model.GradeAssessment, model.GradeProject:
			sa.dailySum += g.Score
			sa.dailyCount++
		case model.GradeUTS:
			sa.uts = g.Score
		case model.GradeUAS:
			sa.uas = g.Score
		}
	}

	round1 := func(v float64) float64 {
		return math.Round(v*10) / 10
	}

	rows := make([]teacherdto.WalasGradeRekapRowDTO, 0, len(students))
	for _, st := range students {
		stName := "Siswa"
		if u, ok := userMap[st.UserID]; ok && u != nil {
			stName = u.Name
		}

		var subjectScores []teacherdto.WalasSubjectScoreDTO
		var finalSum float64
		var finalCount int

		sm := studentGradeMatrix[st.ID]
		for _, sid := range subjectIDs {
			subName := "Mata Pelajaran"
			if sub, ok := subjectMap[sid]; ok && sub != nil {
				subName = sub.Name
			}

			var dAvg, uts, uas, fScore float64
			if sm != nil {
				if sa, ok := sm[sid]; ok && sa != nil {
					if sa.dailyCount > 0 {
						dAvg = round1(sa.dailySum / float64(sa.dailyCount))
					}
					uts = sa.uts
					uas = sa.uas
					fScore = round1((dAvg * 0.6) + (uts * 0.2) + (uas * 0.2))
					finalSum += fScore
					finalCount++
				}
			}

			subjectScores = append(subjectScores, teacherdto.WalasSubjectScoreDTO{
				SubjectID:   sid.Hex(),
				SubjectName: subName,
				DailyAvg:    dAvg,
				UTS:         uts,
				UAS:         uas,
				FinalScore:  fScore,
			})
		}

		overallAvg := float64(0)
		if finalCount > 0 {
			overallAvg = round1(finalSum / float64(finalCount))
		}

		rows = append(rows, teacherdto.WalasGradeRekapRowDTO{
			StudentID:      st.ID.Hex(),
			StudentName:    stName,
			StudentNIS:     st.NIS,
			SubjectScores:  subjectScores,
			OverallAverage: overallAvg,
		})
	}

	// Sort rows by OverallAverage descending to set rank
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].OverallAverage > rows[j].OverallAverage
	})

	for i := range rows {
		rows[i].Rank = i + 1
	}

	return rows, nil
}
