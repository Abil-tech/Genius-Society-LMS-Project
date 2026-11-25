package teachersvc

import (
	"context"
	"sort"
	"time"

	teacherdto "github.com/geniussociety/backend/internal/dto/teacher"
	"github.com/geniussociety/backend/internal/middleware"
	teacherrepo "github.com/geniussociety/backend/internal/repository/teacher"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// DashboardService coordinates data retrieval for the teacher dashboard.
type DashboardService interface {
	GetDashboardData(ctx context.Context, scope *middleware.TeacherScope) (*teacherdto.DashboardResponse, error)
}

type dashboardServiceImpl struct {
	repo      teacherrepo.DashboardRepository
	classRepo teacherrepo.ClassRepository
}

// NewDashboardService creates a new DashboardService instance.
func NewDashboardService(repo teacherrepo.DashboardRepository, classRepo teacherrepo.ClassRepository) DashboardService {
	return &dashboardServiceImpl{
		repo:      repo,
		classRepo: classRepo,
	}
}

func (s *dashboardServiceImpl) GetDashboardData(ctx context.Context, scope *middleware.TeacherScope) (*teacherdto.DashboardResponse, error) {
	teachingClassIDs := scope.GetTeachingClassIDs()

	// 1. KPI Aggregation
	classesCount := len(teachingClassIDs)

	studentsCount, err := s.repo.CountStudentsByClassIDs(ctx, teachingClassIDs)
	if err != nil {
		studentsCount = 0
	}

	activeAssignments, _ := s.repo.CountActiveAssignments(ctx, scope.TeacherID)
	activeAssessments, _ := s.repo.CountActiveAssessments(ctx, scope.TeacherID)
	activeProjects, _ := s.repo.CountActiveProjects(ctx, scope.TeacherID)

	ungradedAssignments, _ := s.repo.CountUngradedAssignmentSubmissions(ctx, scope.TeacherID)
	ungradedProjects, _ := s.repo.CountUngradedProjectSubmissions(ctx, scope.TeacherID)
	ungradedAssessments, _ := s.repo.CountUngradedAssessmentAttempts(ctx, scope.TeacherID)
	pendingGradesCount := ungradedAssignments + ungradedProjects + ungradedAssessments

	kpi := teacherdto.DashboardKPI{
		ClassesCount:           classesCount,
		StudentsCount:          studentsCount,
		PendingGradesCount:     pendingGradesCount,
		ActiveAssignmentsCount: activeAssignments,
		ActiveAssessmentsCount: activeAssessments,
		ActiveProjectsCount:    activeProjects,
	}

	// 2. Today's Schedule
	// Convert Go weekday (0=Sun, 1=Mon, ..., 6=Sat) to model standard (1=Monday ... 7=Sunday)
	goWeekday := time.Now().Weekday()
	dayOfWeek := int(goWeekday)
	if dayOfWeek == 0 {
		dayOfWeek = 7 // Sunday
	}

	schedules, _ := s.repo.FindTodaySchedules(ctx, scope.TeacherID, scope.AcademicYearID, dayOfWeek)
	var todaySchedule []teacherdto.DashboardScheduleItem

	if len(schedules) > 0 {
		var subjectIDs []bson.ObjectID
		var classIDs []bson.ObjectID
		for _, sc := range schedules {
			subjectIDs = append(subjectIDs, sc.SubjectID)
			classIDs = append(classIDs, sc.ClassID)
		}

		subjectMap, _ := s.repo.FindSubjectsByIDs(ctx, subjectIDs)
		classMap, _ := s.repo.FindClassesByIDs(ctx, classIDs)

		for _, sc := range schedules {
			subName := "Mata Pelajaran"
			if sub, ok := subjectMap[sc.SubjectID]; ok && sub != nil {
				subName = sub.Name
			}
			clsName := "Kelas"
			if cls, ok := classMap[sc.ClassID]; ok && cls != nil {
				clsName = cls.Name
			}

			todaySchedule = append(todaySchedule, teacherdto.DashboardScheduleItem{
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
	}
	if todaySchedule == nil {
		todaySchedule = []teacherdto.DashboardScheduleItem{}
	}

	// 3. Recent Submissions (combine assignments & projects, take latest 5-10)
	recentAssignments, _ := s.repo.FindRecentAssignmentSubmissions(ctx, scope.TeacherID, 5)
	recentProjects, _ := s.repo.FindRecentProjectSubmissions(ctx, scope.TeacherID, 5)

	var allSubmissions []teacherdto.DashboardSubmissionItem
	var studentIDs []bson.ObjectID
	var assignmentIDs []bson.ObjectID
	var projectIDs []bson.ObjectID

	for _, as := range recentAssignments {
		studentIDs = append(studentIDs, as.StudentID)
		assignmentIDs = append(assignmentIDs, as.AssignmentID)
	}
	for _, ps := range recentProjects {
		studentIDs = append(studentIDs, ps.StudentID)
		projectIDs = append(projectIDs, ps.ProjectID)
	}

	studentMap, _ := s.repo.FindStudentsByIDs(ctx, studentIDs)
	assignmentMap, _ := s.repo.FindAssignmentsByIDs(ctx, assignmentIDs)
	projectMap, _ := s.repo.FindProjectsByIDs(ctx, projectIDs)

	// Also look up users for student names
	var userIDs []bson.ObjectID
	for _, st := range studentMap {
		if !st.UserID.IsZero() {
			userIDs = append(userIDs, st.UserID)
		}
	}
	userMap, _ := s.classRepo.FindUsersByIDs(ctx, userIDs)

	for _, as := range recentAssignments {
		stName := "Siswa"
		stNIS := ""
		if st, ok := studentMap[as.StudentID]; ok && st != nil {
			stNIS = st.NIS
			if u, uok := userMap[st.UserID]; uok && u != nil {
				stName = u.Name
			}
		}

		title := "Tugas"
		if a, aok := assignmentMap[as.AssignmentID]; aok && a != nil {
			title = a.Title
		}

		allSubmissions = append(allSubmissions, teacherdto.DashboardSubmissionItem{
			ID:           as.ID.Hex(),
			StudentID:    as.StudentID.Hex(),
			StudentName:  stName,
			StudentNIS:   stNIS,
			Type:         "assignment",
			Title:        title,
			SubmittedAt:  as.SubmittedAt.Format(time.RFC3339),
			Graded:       as.Score != nil,
			Score:        as.Score,
			AssignmentID: as.AssignmentID.Hex(),
		})
	}

	for _, ps := range recentProjects {
		stName := "Siswa"
		stNIS := ""
		if st, ok := studentMap[ps.StudentID]; ok && st != nil {
			stNIS = st.NIS
			if u, uok := userMap[st.UserID]; uok && u != nil {
				stName = u.Name
			}
		}

		title := "Projek"
		if p, pok := projectMap[ps.ProjectID]; pok && p != nil {
			title = p.Title
		}

		allSubmissions = append(allSubmissions, teacherdto.DashboardSubmissionItem{
			ID:          ps.ID.Hex(),
			StudentID:   ps.StudentID.Hex(),
			StudentName: stName,
			StudentNIS:  stNIS,
			Type:        "project",
			Title:       title,
			SubmittedAt: ps.SubmittedAt.Format(time.RFC3339),
			Graded:      ps.Score != nil,
			Score:       ps.Score,
			ProjectID:   ps.ProjectID.Hex(),
		})
	}

	// Sort recent submissions by submittedAt descending
	sort.Slice(allSubmissions, func(i, j int) bool {
		return allSubmissions[i].SubmittedAt > allSubmissions[j].SubmittedAt
	})
	if len(allSubmissions) > 6 {
		allSubmissions = allSubmissions[:6]
	}
	if allSubmissions == nil {
		allSubmissions = []teacherdto.DashboardSubmissionItem{}
	}

	// 4. Walas Class Summary (if isWalas)
	var walasClass *teacherdto.DashboardWalasClass
	if scope.IsWalas && len(scope.WalasClassIDs) > 0 {
		wClassID := scope.WalasClassIDs[0]
		cl, err := s.repo.FindClassByID(ctx, wClassID)
		if err == nil && cl != nil {
			stCount, _ := s.repo.CountStudentsByClassID(ctx, wClassID)
			walasClass = &teacherdto.DashboardWalasClass{
				ID:           cl.ID.Hex(),
				Name:         cl.Name,
				Level:        cl.EducationLevel,
				StudentCount: stCount,
			}
		}
	}

	return &teacherdto.DashboardResponse{
		KPI:               kpi,
		TodaySchedule:     todaySchedule,
		RecentSubmissions: allSubmissions,
		IsWalas:           scope.IsWalas,
		WalasClass:        walasClass,
	}, nil
}
