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

func (m *mockClassRepoForService) ExistsByCode(ctx context.Context, academicYearID bson.ObjectID, code string) (bool, error) {
	for _, c := range m.classes {
		if c.AcademicYearID == academicYearID && c.Code == code && c.Status != "archived" {
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

func TestClassService_CRUD(t *testing.T) {
	repo := newMockClassRepoForService()
	svc := NewClassService(repo)
	ctx := context.Background()

	academicYearID := bson.NewObjectID()

	// 1. Create Class
	classA, err := svc.CreateClass(ctx, academicYearID, &dto.ClassCreateDTO{
		Code:           "10-IPA-1",
		Name:           "10 IPA 1",
		EducationLevel: "SMA",
		Capacity:       30,
		Status:         "active",
	})
	if err != nil {
		t.Fatalf("unexpected error creating class: %v", err)
	}

	// 2. Class code duplicate in same academic year should fail
	_, err = svc.CreateClass(ctx, academicYearID, &dto.ClassCreateDTO{
		Code:           "10-IPA-1",
		Name:           "10 IPA 1 Duplicate",
		EducationLevel: "SMA",
		Capacity:       30,
		Status:         "active",
	})
	if err == nil {
		t.Fatal("expected error creating class with duplicate code, got nil")
	}

	// 3. GetClassDetail MUST SUCCEED
	detailA, err := svc.GetClassDetail(ctx, classA.ID)
	if err != nil {
		t.Fatalf("expected GetClassDetail to succeed: %v", err)
	}
	if detailA.ID != classA.ID {
		t.Fatalf("expected class ID %s, got %s", classA.ID.Hex(), detailA.ID.Hex())
	}

	// 4. UpdateClass MUST SUCCEED
	newName := "Updated Name"
	_, err = svc.UpdateClass(ctx, classA.ID, &dto.ClassUpdateDTO{
		Name: &newName,
	})
	if err != nil {
		t.Fatalf("expected UpdateClass to succeed: %v", err)
	}

	// 5. ArchiveClass MUST SUCCEED
	err = svc.ArchiveClass(ctx, classA.ID)
	if err != nil {
		t.Fatalf("expected ArchiveClass to succeed: %v", err)
	}
	if classA.Status != "archived" {
		t.Fatalf("expected class status to be archived, got %s", classA.Status)
	}

	// 6. ListClasses
	_, err = svc.ListClasses(ctx, &dto.ClassListFilter{})
	if err != nil {
		t.Fatalf("unexpected error on ListClasses: %v", err)
	}
}
