package teachersvc

import (
	"context"

	teacherdto "github.com/geniussociety/backend/internal/dto/teacher"
	"github.com/geniussociety/backend/internal/middleware"
	teacherrepo "github.com/geniussociety/backend/internal/repository/teacher"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// ScheduleService defines operations for teacher schedules.
type ScheduleService interface {
	GetTeacherSchedule(ctx context.Context, scope *middleware.TeacherScope) ([]teacherdto.TeacherScheduleItemDTO, error)
}

type scheduleServiceImpl struct {
	repo teacherrepo.ScheduleRepository
}

// NewScheduleService constructs a new ScheduleService.
func NewScheduleService(repo teacherrepo.ScheduleRepository) ScheduleService {
	return &scheduleServiceImpl{repo: repo}
}

func (s *scheduleServiceImpl) GetTeacherSchedule(ctx context.Context, scope *middleware.TeacherScope) ([]teacherdto.TeacherScheduleItemDTO, error) {
	schedules, err := s.repo.FindTeacherSchedules(ctx, scope.TeacherID, scope.AcademicYearID)
	if err != nil {
		return nil, err
	}

	if len(schedules) == 0 {
		return []teacherdto.TeacherScheduleItemDTO{}, nil
	}

	var subjectIDs []bson.ObjectID
	var classIDs []bson.ObjectID
	for _, sc := range schedules {
		subjectIDs = append(subjectIDs, sc.SubjectID)
		classIDs = append(classIDs, sc.ClassID)
	}

	subMap, _ := s.repo.FindSubjectsByIDs(ctx, subjectIDs)
	clsMap, _ := s.repo.FindClassesByIDs(ctx, classIDs)

	results := make([]teacherdto.TeacherScheduleItemDTO, 0, len(schedules))
	for _, sc := range schedules {
		subName := "Mata Pelajaran"
		if sub, ok := subMap[sc.SubjectID]; ok && sub != nil {
			subName = sub.Name
		}
		clsName := "Kelas"
		if cls, ok := clsMap[sc.ClassID]; ok && cls != nil {
			clsName = cls.Name
		}

		results = append(results, teacherdto.TeacherScheduleItemDTO{
			ID:        sc.ID.Hex(),
			Day:       sc.Day,
			StartTime: sc.StartTime,
			EndTime:   sc.EndTime,
			Subject:   subName,
			ClassName: clsName,
			Room:      sc.Room,
			ClassID:   sc.ClassID.Hex(),
			SubjectID: sc.SubjectID.Hex(),
		})
	}

	return results, nil
}
