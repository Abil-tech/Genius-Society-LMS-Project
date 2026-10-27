package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/geniussociety/backend/internal/dto"
	"github.com/geniussociety/backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type mockClassRepoForService struct {
	classes            map[string]*model.Class
	createErr          error
	updateErr          error
	archiveErr         error
	findErr            error
	listErr            error
	lastFilterReceived bson.M
}

func newMockClassRepoForService() *mockClassRepoForService {
	return &mockClassRepoForService{
		classes: make(map[string]*model.Class),
	}
}

func (m *mockClassRepoForService) Create(ctx context.Context, c *model.Class) error {
	if m.createErr != nil {
		return m.createErr
	}
	if c.ID.IsZero() {
		c.ID = bson.NewObjectID()
	}
	m.classes[c.ID.Hex()] = c
	return nil
}

func (m *mockClassRepoForService) ExistsByCode(ctx context.Context, schoolID, academicYearID bson.ObjectID, code string) (bool, error) {
	for _, c := range m.classes {
		if c.SchoolID == schoolID && c.AcademicYearID == academicYearID && c.Code == code && c.Status != "archived" {
			return true, nil
		}
	}
	return false, nil
}

func (m *mockClassRepoForService) FindByID(ctx context.Context, id bson.ObjectID) (*model.Class, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	c, ok := m.classes[id.Hex()]
	if !ok {
		return nil, nil
	}
	return c, nil
}

func (m *mockClassRepoForService) Update(ctx context.Context, id bson.ObjectID, update bson.M) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	c, ok := m.classes[id.Hex()]
	if !ok {
		return errors.New("class not found")
	}
	if name, ok := update["name"].(string); ok {
		c.Name = name
	}
	if status, ok := update["status"].(string); ok {
		c.Status = status
	}
	c.UpdatedAt = time.Now()
	return nil
}

func (m *mockClassRepoForService) Archive(ctx context.Context, id bson.ObjectID) error {
	if m.archiveErr != nil {
		return m.archiveErr
	}
	c, ok := m.classes[id.Hex()]
	if !ok {
		return errors.New("class not found")
	}
	c.Status = "archived"
	c.UpdatedAt = time.Now()
	return nil
}

func (m *mockClassRepoForService) List(ctx context.Context, filter bson.M, page, limit int64) ([]*model.Class, int64, error) {
	m.lastFilterReceived = filter
	if m.listErr != nil {
		return nil, 0, m.listErr
	}
	var res []*model.Class
	schoolID, _ := filter["schoolId"].(bson.ObjectID)
	for _, c := range m.classes {
		if !schoolID.IsZero() && c.SchoolID != schoolID {
			continue
		}
		res = append(res, c)
	}
	return res, int64(len(res)), nil
}

func (m *mockClassRepoForService) CountByFilter(ctx context.Context, filter bson.M) (int64, error) {
	list, total, err := m.List(ctx, filter, 1, 100)
	if err != nil {
		return 0, err
	}
	_ = list
	return total, nil
}

func TestClassService_SchoolIsolation(t *testing.T) {
	repo := newMockClassRepoForService()
	svc := NewClassService(repo)
	ctx := context.Background()

	schoolA := bson.NewObjectID()
	schoolB := bson.NewObjectID()
	academicYearID := bson.NewObjectID()

	// 1. Create Class in School A
	classA, err := svc.CreateClass(ctx, schoolA, academicYearID, &dto.ClassCreateDTO{
		Code:           "10-IPA-1",
		Name:           "10 IPA 1",
		EducationLevel: "SMA",
		Capacity:       30,
		Status:         "active",
	})
	if err != nil {
		t.Fatalf("unexpected error creating class in school A: %v", err)
	}
	if classA.SchoolID != schoolA {
		t.Fatalf("expected classA schoolID to be %s, got %s", schoolA.Hex(), classA.SchoolID.Hex())
	}

	// 2. Class code duplicate in SAME school should fail
	_, err = svc.CreateClass(ctx, schoolA, academicYearID, &dto.ClassCreateDTO{
		Code:           "10-IPA-1",
		Name:           "10 IPA 1 Duplicate",
		EducationLevel: "SMA",
		Capacity:       30,
		Status:         "active",
	})
	if err == nil {
		t.Fatal("expected error creating class with duplicate code in the same school, got nil")
	}

	// 3. Same class code in DIFFERENT school (School B) MUST SUCCEED (school isolation)
	classB, err := svc.CreateClass(ctx, schoolB, academicYearID, &dto.ClassCreateDTO{
		Code:           "10-IPA-1",
		Name:           "10 IPA 1 School B",
		EducationLevel: "SMA",
		Capacity:       32,
		Status:         "active",
	})
	if err != nil {
		t.Fatalf("unexpected error creating same class code in school B: %v", err)
	}
	if classB.SchoolID != schoolB {
		t.Fatalf("expected classB schoolID to be %s, got %s", schoolB.Hex(), classB.SchoolID.Hex())
	}

	// 4. Cross-School GetClassDetail: School B admin trying to view School A class MUST FAIL
	_, err = svc.GetClassDetail(ctx, schoolB, classA.ID)
	if err == nil {
		t.Fatal("expected cross-school GetClassDetail to fail, but it succeeded")
	}

	// Same-School GetClassDetail: School A admin viewing School A class MUST SUCCEED
	detailA, err := svc.GetClassDetail(ctx, schoolA, classA.ID)
	if err != nil {
		t.Fatalf("expected GetClassDetail to succeed for owner school: %v", err)
	}
	if detailA.ID != classA.ID {
		t.Fatalf("expected class ID %s, got %s", classA.ID.Hex(), detailA.ID.Hex())
	}

	// 5. Cross-School UpdateClass: School B admin trying to update School A class MUST FAIL
	newName := "Hacked Name"
	_, err = svc.UpdateClass(ctx, schoolB, classA.ID, &dto.ClassUpdateDTO{
		Name: &newName,
	})
	if err == nil {
		t.Fatal("expected cross-school UpdateClass to fail, but it succeeded")
	}

	// 6. Cross-School ArchiveClass: School B admin trying to archive School A class MUST FAIL
	err = svc.ArchiveClass(ctx, schoolB, classA.ID)
	if err == nil {
		t.Fatal("expected cross-school ArchiveClass to fail, but it succeeded")
	}

	// Same-School ArchiveClass: School A admin archiving School A class MUST SUCCEED
	err = svc.ArchiveClass(ctx, schoolA, classA.ID)
	if err != nil {
		t.Fatalf("expected ArchiveClass to succeed for owner school: %v", err)
	}
	if classA.Status != "archived" {
		t.Fatalf("expected class status to be archived, got %s", classA.Status)
	}

	// 7. ListClasses must enforce school filter
	listRes, err := svc.ListClasses(ctx, schoolA, &dto.ClassListFilter{})
	if err != nil {
		t.Fatalf("unexpected error on ListClasses: %v", err)
	}
	for _, item := range listRes.Items {
		if item.SchoolID != schoolA {
			t.Fatalf("leak detected in ListClasses: returned class belonging to school %s when queried for %s", item.SchoolID.Hex(), schoolA.Hex())
		}
	}
	if repo.lastFilterReceived["schoolId"] != schoolA {
		t.Fatalf("repo did not receive expected schoolID filter: got %v", repo.lastFilterReceived["schoolId"])
	}
}
