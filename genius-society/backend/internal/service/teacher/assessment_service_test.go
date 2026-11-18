package teachersvc

import (
	"context"
	"errors"
	"testing"
	"time"

	teacherdto "github.com/geniussociety/backend/internal/dto/teacher"
	"github.com/geniussociety/backend/internal/handler"
	"github.com/geniussociety/backend/internal/middleware"
	"github.com/geniussociety/backend/internal/model"
	"github.com/geniussociety/backend/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type mockAssessmentRepo struct {
	findAssessmentsFunc func(ctx context.Context, schoolID, teacherID bson.ObjectID, classID, subjectID *bson.ObjectID, search string, page, limit int64) ([]*model.Assessment, int64, error)
	findByIDFunc        func(ctx context.Context, schoolID, teacherID, id bson.ObjectID) (*model.Assessment, error)
	createFunc          func(ctx context.Context, a *model.Assessment) error
	updateFunc          func(ctx context.Context, a *model.Assessment) error
	deleteFunc          func(ctx context.Context, schoolID, teacherID, id bson.ObjectID) error
	countAttemptsFunc   func(ctx context.Context, schoolID, assessmentID bson.ObjectID) (int64, error)

	findQuestionsFunc    func(ctx context.Context, assessmentID bson.ObjectID) ([]*model.AssessmentQuestion, error)
	findQuestionByIDFunc func(ctx context.Context, assessmentID, questionID bson.ObjectID) (*model.AssessmentQuestion, error)
	createQuestionFunc   func(ctx context.Context, q *model.AssessmentQuestion) error
	updateQuestionFunc   func(ctx context.Context, q *model.AssessmentQuestion) error
	deleteQuestionFunc   func(ctx context.Context, assessmentID, questionID bson.ObjectID) error
	reorderQuestionsFunc func(ctx context.Context, assessmentID bson.ObjectID, orders map[bson.ObjectID]int) error

	findAttemptsFunc    func(ctx context.Context, schoolID, assessmentID bson.ObjectID) ([]*model.AssessmentAttempt, error)
	findAttemptByIDFunc func(ctx context.Context, schoolID, attemptID bson.ObjectID) (*model.AssessmentAttempt, error)
	updateAttemptFunc   func(ctx context.Context, attempt *model.AssessmentAttempt) error

	findAnswersFunc    func(ctx context.Context, attemptID bson.ObjectID) ([]*model.AssessmentAnswer, error)
	findAnswerByIDFunc func(ctx context.Context, answerID bson.ObjectID) (*model.AssessmentAnswer, error)

	findStudentsFunc    func(ctx context.Context, schoolID bson.ObjectID, classIDs []bson.ObjectID) ([]*model.Student, error)
	findStudentByIDFunc func(ctx context.Context, schoolID, studentID bson.ObjectID) (*model.Student, error)
	findUsersMapFunc    func(ctx context.Context, userIDs []bson.ObjectID) (map[bson.ObjectID]*model.User, error)
	findClassesMapFunc  func(ctx context.Context, schoolID bson.ObjectID, classIDs []bson.ObjectID) (map[bson.ObjectID]*model.Class, error)
	findSubjectsMapFunc func(ctx context.Context, subjectIDs []bson.ObjectID) (map[bson.ObjectID]*model.Subject, error)

	gradeEssayInTxFunc func(ctx context.Context, schoolID, teacherID bson.ObjectID, answer *model.AssessmentAnswer, attempt *model.AssessmentAttempt, grade *model.Grade, audit *model.AuditLog) error
	regradeMCQInTxFunc func(ctx context.Context, schoolID, teacherID bson.ObjectID, question *model.AssessmentQuestion, attempts []*model.AssessmentAttempt, answers []*model.AssessmentAnswer, grades []*model.Grade, audit *model.AuditLog) error
}

func (m *mockAssessmentRepo) FindAssessments(ctx context.Context, schoolID, teacherID bson.ObjectID, classID, subjectID *bson.ObjectID, search string, page, limit int64) ([]*model.Assessment, int64, error) {
	if m.findAssessmentsFunc != nil {
		return m.findAssessmentsFunc(ctx, schoolID, teacherID, classID, subjectID, search, page, limit)
	}
	return nil, 0, nil
}

func (m *mockAssessmentRepo) FindAssessmentByID(ctx context.Context, schoolID, teacherID, id bson.ObjectID) (*model.Assessment, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, schoolID, teacherID, id)
	}
	return nil, nil
}

func (m *mockAssessmentRepo) CreateAssessment(ctx context.Context, a *model.Assessment) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, a)
	}
	return nil
}

func (m *mockAssessmentRepo) UpdateAssessment(ctx context.Context, a *model.Assessment) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, a)
	}
	return nil
}

func (m *mockAssessmentRepo) DeleteAssessment(ctx context.Context, schoolID, teacherID, id bson.ObjectID) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, schoolID, teacherID, id)
	}
	return nil
}

func (m *mockAssessmentRepo) CountAttemptsByAssessmentID(ctx context.Context, schoolID, assessmentID bson.ObjectID) (int64, error) {
	if m.countAttemptsFunc != nil {
		return m.countAttemptsFunc(ctx, schoolID, assessmentID)
	}
	return 0, nil
}

func (m *mockAssessmentRepo) FindQuestions(ctx context.Context, assessmentID bson.ObjectID) ([]*model.AssessmentQuestion, error) {
	if m.findQuestionsFunc != nil {
		return m.findQuestionsFunc(ctx, assessmentID)
	}
	return nil, nil
}

func (m *mockAssessmentRepo) FindQuestionByID(ctx context.Context, assessmentID, questionID bson.ObjectID) (*model.AssessmentQuestion, error) {
	if m.findQuestionByIDFunc != nil {
		return m.findQuestionByIDFunc(ctx, assessmentID, questionID)
	}
	return nil, nil
}

func (m *mockAssessmentRepo) CreateQuestion(ctx context.Context, q *model.AssessmentQuestion) error {
	if m.createQuestionFunc != nil {
		return m.createQuestionFunc(ctx, q)
	}
	return nil
}

func (m *mockAssessmentRepo) UpdateQuestion(ctx context.Context, q *model.AssessmentQuestion) error {
	if m.updateQuestionFunc != nil {
		return m.updateQuestionFunc(ctx, q)
	}
	return nil
}

func (m *mockAssessmentRepo) DeleteQuestion(ctx context.Context, assessmentID, questionID bson.ObjectID) error {
	if m.deleteQuestionFunc != nil {
		return m.deleteQuestionFunc(ctx, assessmentID, questionID)
	}
	return nil
}

func (m *mockAssessmentRepo) ReorderQuestions(ctx context.Context, assessmentID bson.ObjectID, orders map[bson.ObjectID]int) error {
	if m.reorderQuestionsFunc != nil {
		return m.reorderQuestionsFunc(ctx, assessmentID, orders)
	}
	return nil
}

func (m *mockAssessmentRepo) FindAttempts(ctx context.Context, schoolID, assessmentID bson.ObjectID) ([]*model.AssessmentAttempt, error) {
	if m.findAttemptsFunc != nil {
		return m.findAttemptsFunc(ctx, schoolID, assessmentID)
	}
	return nil, nil
}

func (m *mockAssessmentRepo) FindAttemptByID(ctx context.Context, schoolID, attemptID bson.ObjectID) (*model.AssessmentAttempt, error) {
	if m.findAttemptByIDFunc != nil {
		return m.findAttemptByIDFunc(ctx, schoolID, attemptID)
	}
	return nil, nil
}

func (m *mockAssessmentRepo) UpdateAttempt(ctx context.Context, attempt *model.AssessmentAttempt) error {
	if m.updateAttemptFunc != nil {
		return m.updateAttemptFunc(ctx, attempt)
	}
	return nil
}

func (m *mockAssessmentRepo) FindAnswersByAttemptID(ctx context.Context, attemptID bson.ObjectID) ([]*model.AssessmentAnswer, error) {
	if m.findAnswersFunc != nil {
		return m.findAnswersFunc(ctx, attemptID)
	}
	return nil, nil
}

func (m *mockAssessmentRepo) FindAnswerByID(ctx context.Context, answerID bson.ObjectID) (*model.AssessmentAnswer, error) {
	if m.findAnswerByIDFunc != nil {
		return m.findAnswerByIDFunc(ctx, answerID)
	}
	return nil, nil
}

func (m *mockAssessmentRepo) FindStudentsByClassIDs(ctx context.Context, schoolID bson.ObjectID, classIDs []bson.ObjectID) ([]*model.Student, error) {
	if m.findStudentsFunc != nil {
		return m.findStudentsFunc(ctx, schoolID, classIDs)
	}
	return nil, nil
}

func (m *mockAssessmentRepo) FindStudentByID(ctx context.Context, schoolID, studentID bson.ObjectID) (*model.Student, error) {
	if m.findStudentByIDFunc != nil {
		return m.findStudentByIDFunc(ctx, schoolID, studentID)
	}
	return nil, nil
}

func (m *mockAssessmentRepo) FindUsersMap(ctx context.Context, userIDs []bson.ObjectID) (map[bson.ObjectID]*model.User, error) {
	if m.findUsersMapFunc != nil {
		return m.findUsersMapFunc(ctx, userIDs)
	}
	return make(map[bson.ObjectID]*model.User), nil
}

func (m *mockAssessmentRepo) FindClassesMap(ctx context.Context, schoolID bson.ObjectID, classIDs []bson.ObjectID) (map[bson.ObjectID]*model.Class, error) {
	if m.findClassesMapFunc != nil {
		return m.findClassesMapFunc(ctx, schoolID, classIDs)
	}
	return make(map[bson.ObjectID]*model.Class), nil
}

func (m *mockAssessmentRepo) FindSubjectsMap(ctx context.Context, subjectIDs []bson.ObjectID) (map[bson.ObjectID]*model.Subject, error) {
	if m.findSubjectsMapFunc != nil {
		return m.findSubjectsMapFunc(ctx, subjectIDs)
	}
	return make(map[bson.ObjectID]*model.Subject), nil
}

func (m *mockAssessmentRepo) GradeEssayInTransaction(ctx context.Context, schoolID, teacherID bson.ObjectID, answer *model.AssessmentAnswer, attempt *model.AssessmentAttempt, grade *model.Grade, audit *model.AuditLog) error {
	if m.gradeEssayInTxFunc != nil {
		return m.gradeEssayInTxFunc(ctx, schoolID, teacherID, answer, attempt, grade, audit)
	}
	return nil
}

func (m *mockAssessmentRepo) RegradeMCQInTransaction(ctx context.Context, schoolID, teacherID bson.ObjectID, question *model.AssessmentQuestion, attempts []*model.AssessmentAttempt, answers []*model.AssessmentAnswer, grades []*model.Grade, audit *model.AuditLog) error {
	if m.regradeMCQInTxFunc != nil {
		return m.regradeMCQInTxFunc(ctx, schoolID, teacherID, question, attempts, answers, grades, audit)
	}
	return nil
}

func TestAssessmentService_Delete_ConflictWhenHasAttempts(t *testing.T) {
	schoolID := bson.NewObjectID()
	teacherID := bson.NewObjectID()
	asID := bson.NewObjectID()

	scope := &middleware.TeacherScope{
		SchoolID:  schoolID,
		TeacherID: teacherID,
	}

	repo := &mockAssessmentRepo{
		findByIDFunc: func(ctx context.Context, sID, tID, id bson.ObjectID) (*model.Assessment, error) {
			return &model.Assessment{ID: asID, SchoolID: schoolID, TeacherID: teacherID}, nil
		},
		countAttemptsFunc: func(ctx context.Context, sID, aID bson.ObjectID) (int64, error) {
			return 1, nil // Has 1 attempt!
		},
	}

	svc := NewAssessmentService(repo, nil, nil)
	err := svc.DeleteAssessment(context.Background(), scope, asID)
	if err == nil {
		t.Fatal("expected 409 conflict when deleting assessment with attempts")
	}

	var dErr *handler.DomainError
	if !errors.As(err, &dErr) || dErr.Status != 409 {
		t.Errorf("expected DomainError with status 409, got: %v", err)
	}
}

func TestAssessmentService_QuestionMutation_ConflictAfterAttempt(t *testing.T) {
	schoolID := bson.NewObjectID()
	teacherID := bson.NewObjectID()
	asID := bson.NewObjectID()

	scope := &middleware.TeacherScope{
		SchoolID:  schoolID,
		TeacherID: teacherID,
	}

	repo := &mockAssessmentRepo{
		findByIDFunc: func(ctx context.Context, sID, tID, id bson.ObjectID) (*model.Assessment, error) {
			return &model.Assessment{ID: asID, SchoolID: schoolID, TeacherID: teacherID}, nil
		},
		countAttemptsFunc: func(ctx context.Context, sID, aID bson.ObjectID) (int64, error) {
			return 2, nil // Has attempts!
		},
	}

	svc := NewAssessmentService(repo, nil, nil)

	// Attempting to add question after attempts exist
	correctOpt := 0
	_, err := svc.AddQuestion(context.Background(), scope, asID, &teacherdto.QuestionCreateUpdatePayload{
		Type:          model.QuestionMCQ,
		Text:          "Soal tambahan",
		Options:       []string{"A", "B"},
		CorrectOption: &correctOpt,
		Weight:        10,
	})
	if err == nil {
		t.Fatal("expected 409 conflict when adding question after attempts")
	}

	var dErr *handler.DomainError
	if !errors.As(err, &dErr) || dErr.Status != 409 {
		t.Errorf("expected DomainError with status 409, got: %v", err)
	}
}

func TestAssessmentService_Publish_Validation(t *testing.T) {
	schoolID := bson.NewObjectID()
	teacherID := bson.NewObjectID()
	asID := bson.NewObjectID()
	classID := bson.NewObjectID()
	subjectID := bson.NewObjectID()

	scope := &middleware.TeacherScope{
		SchoolID:  schoolID,
		TeacherID: teacherID,
		TeachingPairs: []repository.ClassSubjectPair{
			{ClassID: classID, SubjectID: subjectID},
		},
	}

	// 1. Publish without questions -> validation error
	repo := &mockAssessmentRepo{
		findByIDFunc: func(ctx context.Context, sID, tID, id bson.ObjectID) (*model.Assessment, error) {
			now := time.Now()
			end := now.Add(2 * time.Hour)
			return &model.Assessment{
				ID:              asID,
				SchoolID:        schoolID,
				TeacherID:       teacherID,
				SubjectID:       subjectID,
				ClassIDs:        []bson.ObjectID{classID},
				DurationMinutes: 60,
				StartAt:         &now,
				EndAt:           &end,
			}, nil
		},
		findQuestionsFunc: func(ctx context.Context, aID bson.ObjectID) ([]*model.AssessmentQuestion, error) {
			return []*model.AssessmentQuestion{}, nil // No questions!
		},
	}

	svc := NewAssessmentService(repo, nil, nil)
	_, err := svc.PublishAssessment(context.Background(), scope, asID)
	if err == nil {
		t.Fatal("expected validation error when publishing assessment with 0 questions")
	}
}

func TestAssessmentService_GradeEssay_ExceedWeight(t *testing.T) {
	schoolID := bson.NewObjectID()
	teacherID := bson.NewObjectID()
	asID := bson.NewObjectID()
	qID := bson.NewObjectID()
	ansID := bson.NewObjectID()
	attID := bson.NewObjectID()

	scope := &middleware.TeacherScope{
		SchoolID:  schoolID,
		TeacherID: teacherID,
	}

	repo := &mockAssessmentRepo{
		findByIDFunc: func(ctx context.Context, sID, tID, id bson.ObjectID) (*model.Assessment, error) {
			return &model.Assessment{ID: asID, SchoolID: schoolID, TeacherID: teacherID}, nil
		},
		findAnswerByIDFunc: func(ctx context.Context, id bson.ObjectID) (*model.AssessmentAnswer, error) {
			return &model.AssessmentAnswer{ID: ansID, AttemptID: attID, QuestionID: qID}, nil
		},
		findAttemptByIDFunc: func(ctx context.Context, sID, id bson.ObjectID) (*model.AssessmentAttempt, error) {
			return &model.AssessmentAttempt{ID: attID, AssessmentID: asID, SchoolID: schoolID}, nil
		},
		findQuestionByIDFunc: func(ctx context.Context, aID, id bson.ObjectID) (*model.AssessmentQuestion, error) {
			return &model.AssessmentQuestion{ID: qID, AssessmentID: asID, Type: model.QuestionEssay, Weight: 20.0}, nil
		},
	}

	svc := NewAssessmentService(repo, nil, nil)

	// Score 25.0 exceeds weight 20.0 -> must fail with validation error!
	err := svc.GradeEssayAnswer(context.Background(), scope, asID, ansID, 25.0, "Bagus")
	if err == nil {
		t.Fatal("expected validation error when score exceeds question weight")
	}
}
