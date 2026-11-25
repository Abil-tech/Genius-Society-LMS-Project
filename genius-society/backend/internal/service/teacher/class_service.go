package teachersvc

import (
	"context"
	"math"

	teacherdto "github.com/geniussociety/backend/internal/dto/teacher"
	"github.com/geniussociety/backend/internal/handler"
	"github.com/geniussociety/backend/internal/middleware"
	teacherrepo "github.com/geniussociety/backend/internal/repository/teacher"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// ClassService defines teacher class operations.
type ClassService interface {
	GetClasses(ctx context.Context, scope *middleware.TeacherScope) ([]teacherdto.TeacherClassDTO, error)
	GetClassDetail(ctx context.Context, scope *middleware.TeacherScope, classID bson.ObjectID) (*teacherdto.TeacherClassDTO, error)
	GetClassStudents(ctx context.Context, scope *middleware.TeacherScope, classID bson.ObjectID, search string, page, limit int64) ([]teacherdto.TeacherStudentDTO, *teacherdto.ClassStudentsMeta, error)
	GetTeacherSubjects(ctx context.Context, scope *middleware.TeacherScope) ([]teacherdto.TeacherSubjectDTO, error)
}

type classServiceImpl struct {
	repo teacherrepo.ClassRepository
}

// NewClassService creates a new ClassService instance.
func NewClassService(repo teacherrepo.ClassRepository) ClassService {
	return &classServiceImpl{repo: repo}
}

func (s *classServiceImpl) GetClasses(ctx context.Context, scope *middleware.TeacherScope) ([]teacherdto.TeacherClassDTO, error) {
	if scope == nil {
		return []teacherdto.TeacherClassDTO{}, nil
	}

	// 1. Gather all accessible class IDs (teaching classes + walas classes)
	classIDMap := make(map[bson.ObjectID]bool)
	for _, cid := range scope.GetTeachingClassIDs() {
		classIDMap[cid] = true
	}
	for _, wid := range scope.WalasClassIDs {
		classIDMap[wid] = true
	}

	var allClassIDs []bson.ObjectID
	for cid := range classIDMap {
		allClassIDs = append(allClassIDs, cid)
	}

	if len(allClassIDs) == 0 {
		return []teacherdto.TeacherClassDTO{}, nil
	}

	// 2. Fetch classes
	classes, err := s.repo.FindClassesByIDs(ctx, scope.SchoolID, allClassIDs)
	if err != nil {
		return nil, err
	}

	// 3. Collect all teaching subject IDs to resolve subject names in one batch
	var allSubjectIDs []bson.ObjectID
	subjectSeen := make(map[bson.ObjectID]bool)
	for _, p := range scope.TeachingPairs {
		if !subjectSeen[p.SubjectID] {
			subjectSeen[p.SubjectID] = true
			allSubjectIDs = append(allSubjectIDs, p.SubjectID)
		}
	}

	subjectMap, err := s.repo.FindSubjectsByIDs(ctx, allSubjectIDs)
	if err != nil {
		return nil, err
	}

	// 4. Build response DTOs
	result := make([]teacherdto.TeacherClassDTO, 0, len(classes))
	for _, cl := range classes {
		studentCount, _ := s.repo.CountStudentsByClassID(ctx, scope.SchoolID, cl.ID)

		// Find subjects taught by this teacher in this specific class
		var subjects []teacherdto.SubjectBrief
		subSeenInClass := make(map[bson.ObjectID]bool)
		for _, p := range scope.TeachingPairs {
			if p.ClassID == cl.ID && !subSeenInClass[p.SubjectID] {
				subSeenInClass[p.SubjectID] = true
				if sub, ok := subjectMap[p.SubjectID]; ok && sub != nil {
					subjects = append(subjects, teacherdto.SubjectBrief{
						ID:   sub.ID.Hex(),
						Code: sub.Code,
						Name: sub.Name,
					})
				}
			}
		}
		if subjects == nil {
			subjects = []teacherdto.SubjectBrief{}
		}

		dept := ""
		if cl.Department != nil {
			dept = *cl.Department
		}

		result = append(result, teacherdto.TeacherClassDTO{
			ID:             cl.ID.Hex(),
			Code:           cl.Code,
			Name:           cl.Name,
			EducationLevel: cl.EducationLevel,
			Department:     dept,
			AcademicYearID: cl.AcademicYearID.Hex(),
			StudentCount:   studentCount,
			Subjects:       subjects,
			IsWalas:        scope.IsWalasOf(cl.ID),
			Status:         cl.Status,
		})
	}

	return result, nil
}

func (s *classServiceImpl) GetClassDetail(ctx context.Context, scope *middleware.TeacherScope, classID bson.ObjectID) (*teacherdto.TeacherClassDTO, error) {
	if scope == nil || !scope.CanAccessClass(classID) {
		return nil, handler.ErrNotFound
	}

	cl, err := s.repo.FindClassByID(ctx, scope.SchoolID, classID)
	if err != nil {
		return nil, err
	}
	if cl == nil {
		return nil, handler.ErrNotFound
	}

	studentCount, _ := s.repo.CountStudentsByClassID(ctx, scope.SchoolID, cl.ID)

	var classSubjectIDs []bson.ObjectID
	for _, p := range scope.TeachingPairs {
		if p.ClassID == cl.ID {
			classSubjectIDs = append(classSubjectIDs, p.SubjectID)
		}
	}

	subjectMap, err := s.repo.FindSubjectsByIDs(ctx, classSubjectIDs)
	if err != nil {
		return nil, err
	}

	var subjects []teacherdto.SubjectBrief
	subSeen := make(map[bson.ObjectID]bool)
	for _, p := range scope.TeachingPairs {
		if p.ClassID == cl.ID && !subSeen[p.SubjectID] {
			subSeen[p.SubjectID] = true
			if sub, ok := subjectMap[p.SubjectID]; ok && sub != nil {
				subjects = append(subjects, teacherdto.SubjectBrief{
					ID:   sub.ID.Hex(),
					Code: sub.Code,
					Name: sub.Name,
				})
			}
		}
	}
	if subjects == nil {
		subjects = []teacherdto.SubjectBrief{}
	}

	dept := ""
	if cl.Department != nil {
		dept = *cl.Department
	}

	return &teacherdto.TeacherClassDTO{
		ID:             cl.ID.Hex(),
		Code:           cl.Code,
		Name:           cl.Name,
		EducationLevel: cl.EducationLevel,
		Department:     dept,
		AcademicYearID: cl.AcademicYearID.Hex(),
		StudentCount:   studentCount,
		Subjects:       subjects,
		IsWalas:        scope.IsWalasOf(cl.ID),
		Status:         cl.Status,
	}, nil
}

func (s *classServiceImpl) GetClassStudents(ctx context.Context, scope *middleware.TeacherScope, classID bson.ObjectID, search string, page, limit int64) ([]teacherdto.TeacherStudentDTO, *teacherdto.ClassStudentsMeta, error) {
	if scope == nil || !scope.CanAccessClass(classID) {
		return nil, nil, handler.ErrNotFound
	}

	cl, err := s.repo.FindClassByID(ctx, scope.SchoolID, classID)
	if err != nil {
		return nil, nil, err
	}
	if cl == nil {
		return nil, nil, handler.ErrNotFound
	}

	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if page <= 0 {
		page = 1
	}

	students, total, err := s.repo.FindStudentsByClassID(ctx, scope.SchoolID, classID, search, page, limit)
	if err != nil {
		return nil, nil, err
	}

	var userIDs []bson.ObjectID
	for _, st := range students {
		if !st.UserID.IsZero() {
			userIDs = append(userIDs, st.UserID)
		}
	}

	userMap, err := s.repo.FindUsersByIDs(ctx, userIDs)
	if err != nil {
		return nil, nil, err
	}

	result := make([]teacherdto.TeacherStudentDTO, 0, len(students))
	for _, st := range students {
		name := ""
		email := ""
		hasAccount := false

		if !st.UserID.IsZero() {
			if u, ok := userMap[st.UserID]; ok && u != nil {
				name = u.Name
				email = u.Email
				hasAccount = true
			}
		}

		result = append(result, teacherdto.TeacherStudentDTO{
			ID:         st.ID.Hex(),
			StudentID:  st.ID.Hex(),
			Name:       name,
			NIS:        st.NIS,
			NISN:       st.NISN,
			Gender:     st.Gender,
			ClassID:    st.ClassID.Hex(),
			ClassName:  cl.Name,
			Email:      email,
			Status:     st.Status,
			HasAccount: hasAccount,
		})
	}

	totalPages := int64(0)
	if total > 0 {
		totalPages = int64(math.Ceil(float64(total) / float64(limit)))
	}

	meta := &teacherdto.ClassStudentsMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	return result, meta, nil
}

func (s *classServiceImpl) GetTeacherSubjects(ctx context.Context, scope *middleware.TeacherScope) ([]teacherdto.TeacherSubjectDTO, error) {
	if scope == nil {
		return []teacherdto.TeacherSubjectDTO{}, nil
	}

	subjectIDs := scope.GetTeachingSubjectIDs()
	if len(subjectIDs) == 0 {
		return []teacherdto.TeacherSubjectDTO{}, nil
	}

	subjectMap, err := s.repo.FindSubjectsByIDs(ctx, subjectIDs)
	if err != nil {
		return nil, err
	}

	result := make([]teacherdto.TeacherSubjectDTO, 0, len(subjectIDs))
	for _, sID := range subjectIDs {
		sub, ok := subjectMap[sID]
		if !ok || sub == nil {
			continue
		}

		// Find class IDs taught for this subject
		var classIDs []string
		classSeen := make(map[bson.ObjectID]bool)
		for _, p := range scope.TeachingPairs {
			if p.SubjectID == sID && !classSeen[p.ClassID] {
				classSeen[p.ClassID] = true
				classIDs = append(classIDs, p.ClassID.Hex())
			}
		}
		if classIDs == nil {
			classIDs = []string{}
		}

		result = append(result, teacherdto.TeacherSubjectDTO{
			ID:       sub.ID.Hex(),
			Code:     sub.Code,
			Name:     sub.Name,
			ClassIDs: classIDs,
		})
	}

	return result, nil
}
