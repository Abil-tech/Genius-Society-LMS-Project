package teachersvc

import (
	"context"
	"errors"
	"testing"

	"github.com/geniussociety/backend/internal/handler"
	"github.com/geniussociety/backend/internal/middleware"
	"github.com/geniussociety/backend/internal/model"
	"github.com/geniussociety/backend/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type mockClassRepo struct {
	classesFunc       func(ctx context.Context, schoolID bson.ObjectID, classIDs []bson.ObjectID) ([]*model.Class, error)
	classByIDFunc     func(ctx context.Context, schoolID, classID bson.ObjectID) (*model.Class, error)
	countStudentsFunc func(ctx context.Context, schoolID, classID bson.ObjectID) (int64, error)
	studentsFunc      func(ctx context.Context, schoolID, classID bson.ObjectID, search string, page, limit int64) ([]*model.Student, int64, error)
	usersFunc         func(ctx context.Context, userIDs []bson.ObjectID) (map[bson.ObjectID]*model.User, error)
	subjectsFunc      func(ctx context.Context, subjectIDs []bson.ObjectID) (map[bson.ObjectID]*model.Subject, error)
}

func (m *mockClassRepo) FindClassesByIDs(ctx context.Context, schoolID bson.ObjectID, classIDs []bson.ObjectID) ([]*model.Class, error) {
	if m.classesFunc != nil {
		return m.classesFunc(ctx, schoolID, classIDs)
	}
	return nil, nil
}

func (m *mockClassRepo) FindClassByID(ctx context.Context, schoolID, classID bson.ObjectID) (*model.Class, error) {
	if m.classByIDFunc != nil {
		return m.classByIDFunc(ctx, schoolID, classID)
	}
	return nil, nil
}

func (m *mockClassRepo) CountStudentsByClassID(ctx context.Context, schoolID, classID bson.ObjectID) (int64, error) {
	if m.countStudentsFunc != nil {
		return m.countStudentsFunc(ctx, schoolID, classID)
	}
	return 0, nil
}

func (m *mockClassRepo) FindStudentsByClassID(ctx context.Context, schoolID, classID bson.ObjectID, search string, page, limit int64) ([]*model.Student, int64, error) {
	if m.studentsFunc != nil {
		return m.studentsFunc(ctx, schoolID, classID, search, page, limit)
	}
	return nil, 0, nil
}

func (m *mockClassRepo) FindUsersByIDs(ctx context.Context, userIDs []bson.ObjectID) (map[bson.ObjectID]*model.User, error) {
	if m.usersFunc != nil {
		return m.usersFunc(ctx, userIDs)
	}
	return make(map[bson.ObjectID]*model.User), nil
}

func (m *mockClassRepo) FindSubjectsByIDs(ctx context.Context, subjectIDs []bson.ObjectID) (map[bson.ObjectID]*model.Subject, error) {
	if m.subjectsFunc != nil {
		return m.subjectsFunc(ctx, subjectIDs)
	}
	return make(map[bson.ObjectID]*model.Subject), nil
}

func TestGetClasses(t *testing.T) {
	schoolID := bson.NewObjectID()
	class1 := bson.NewObjectID()
	classWalas := bson.NewObjectID()
	subj1 := bson.NewObjectID()

	scope := &middleware.TeacherScope{
		SchoolID: schoolID,
		TeachingPairs: []repository.ClassSubjectPair{
			{ClassID: class1, SubjectID: subj1},
		},
		WalasClassIDs: []bson.ObjectID{classWalas},
		IsWalas:       true,
	}

	repo := &mockClassRepo{
		classesFunc: func(ctx context.Context, sID bson.ObjectID, cIDs []bson.ObjectID) ([]*model.Class, error) {
			if len(cIDs) != 2 {
				t.Errorf("expected 2 class IDs, got %d", len(cIDs))
			}
			return []*model.Class{
				{ID: class1, Name: "X-IPA-1", Code: "X-IPA-1", EducationLevel: "SMA"},
				{ID: classWalas, Name: "X-IPA-2", Code: "X-IPA-2", EducationLevel: "SMA"},
			}, nil
		},
		subjectsFunc: func(ctx context.Context, sIDs []bson.ObjectID) (map[bson.ObjectID]*model.Subject, error) {
			return map[bson.ObjectID]*model.Subject{
				subj1: {ID: subj1, Code: "MTK", Name: "Matematika"},
			}, nil
		},
		countStudentsFunc: func(ctx context.Context, sID, cID bson.ObjectID) (int64, error) {
			if cID == class1 {
				return 30, nil
			}
			return 25, nil
		},
	}

	svc := NewClassService(repo)
	res, err := svc.GetClasses(context.Background(), scope)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 2 {
		t.Fatalf("expected 2 classes, got %d", len(res))
	}

	// Verify class1 has subjects and not walas
	if res[0].ID == class1.Hex() {
		if !res[0].IsWalas && res[0].StudentCount == 30 && len(res[0].Subjects) == 1 {
			// correct
		} else {
			t.Errorf("unexpected class1 dto: %+v", res[0])
		}
	}

	// Verify classWalas is marked walas
	var walasFound bool
	for _, c := range res {
		if c.ID == classWalas.Hex() && c.IsWalas {
			walasFound = true
		}
	}
	if !walasFound {
		t.Errorf("expected walas class to have isWalas=true")
	}
}

func TestGetClassDetail_AccessControl(t *testing.T) {
	schoolID := bson.NewObjectID()
	classAllowed := bson.NewObjectID()
	classForbidden := bson.NewObjectID()

	scope := &middleware.TeacherScope{
		SchoolID: schoolID,
		TeachingPairs: []repository.ClassSubjectPair{
			{ClassID: classAllowed, SubjectID: bson.NewObjectID()},
		},
	}

	repo := &mockClassRepo{
		classByIDFunc: func(ctx context.Context, sID, cID bson.ObjectID) (*model.Class, error) {
			if cID == classAllowed {
				return &model.Class{ID: classAllowed, Name: "Allowed"}, nil
			}
			return nil, nil
		},
	}

	svc := NewClassService(repo)

	// Allowed class
	detail, err := svc.GetClassDetail(context.Background(), scope, classAllowed)
	if err != nil || detail == nil {
		t.Errorf("expected success for allowed class, got err: %v", err)
	}

	// Forbidden class -> 404
	_, err = svc.GetClassDetail(context.Background(), scope, classForbidden)
	if !errors.Is(err, handler.ErrNotFound) {
		t.Errorf("expected ErrNotFound for forbidden class, got: %v", err)
	}
}

func TestGetClassStudents_BulkUserLookupAndAccountFlag(t *testing.T) {
	schoolID := bson.NewObjectID()
	classID := bson.NewObjectID()
	user1 := bson.NewObjectID()

	scope := &middleware.TeacherScope{
		SchoolID: schoolID,
		TeachingPairs: []repository.ClassSubjectPair{
			{ClassID: classID, SubjectID: bson.NewObjectID()},
		},
	}

	studentWithAccount := &model.Student{
		ID:       bson.NewObjectID(),
		ClassID:  classID,
		UserID:   user1,
		NIS:      "1001",
		NISN:     "001001",
		Status:   "active",
		SchoolID: schoolID,
	}
	studentWithoutAccount := &model.Student{
		ID:       bson.NewObjectID(),
		ClassID:  classID,
		UserID:   bson.NilObjectID, // no account!
		NIS:      "1002",
		Status:   "active",
		SchoolID: schoolID,
	}

	var bulkLookupCount int
	repo := &mockClassRepo{
		classByIDFunc: func(ctx context.Context, sID, cID bson.ObjectID) (*model.Class, error) {
			return &model.Class{ID: classID, Name: "X-IPA-1"}, nil
		},
		studentsFunc: func(ctx context.Context, sID, cID bson.ObjectID, search string, page, limit int64) ([]*model.Student, int64, error) {
			return []*model.Student{studentWithAccount, studentWithoutAccount}, 2, nil
		},
		usersFunc: func(ctx context.Context, uIDs []bson.ObjectID) (map[bson.ObjectID]*model.User, error) {
			bulkLookupCount++
			if len(uIDs) != 1 || uIDs[0] != user1 {
				t.Errorf("expected bulk lookup for user1 only, got %v", uIDs)
			}
			return map[bson.ObjectID]*model.User{
				user1: {ID: user1, Name: "Ahmad Siswa", Email: "ahmad@test.com"},
			}, nil
		},
	}

	svc := NewClassService(repo)
	students, meta, err := svc.GetClassStudents(context.Background(), scope, classID, "", 1, 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bulkLookupCount != 1 {
		t.Errorf("expected exactly 1 bulk user lookup, got %d", bulkLookupCount)
	}
	if len(students) != 2 {
		t.Fatalf("expected 2 students, got %d", len(students))
	}
	if meta.Total != 2 || meta.TotalPages != 1 {
		t.Errorf("unexpected meta: %+v", meta)
	}

	// Student 1: has account
	if students[0].Name != "Ahmad Siswa" || !students[0].HasAccount || students[0].Email != "ahmad@test.com" {
		t.Errorf("unexpected studentWithAccount dto: %+v", students[0])
	}

	// Student 2: no account -> name must be empty string, hasAccount=false
	if students[1].Name != "" || students[1].HasAccount || students[1].Email != "" {
		t.Errorf("unexpected studentWithoutAccount dto (must have empty name and false hasAccount): %+v", students[1])
	}
}
