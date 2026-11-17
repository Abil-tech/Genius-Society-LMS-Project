package teacher

import (
	"context"
	"regexp"

	"github.com/geniussociety/backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// ClassRepository defines read operations needed for teacher class views.
type ClassRepository interface {
	FindClassesByIDs(ctx context.Context, schoolID bson.ObjectID, classIDs []bson.ObjectID) ([]*model.Class, error)
	FindClassByID(ctx context.Context, schoolID, classID bson.ObjectID) (*model.Class, error)
	CountStudentsByClassID(ctx context.Context, schoolID, classID bson.ObjectID) (int64, error)
	FindStudentsByClassID(ctx context.Context, schoolID, classID bson.ObjectID, search string, page, limit int64) ([]*model.Student, int64, error)
	FindUsersByIDs(ctx context.Context, userIDs []bson.ObjectID) (map[bson.ObjectID]*model.User, error)
	FindSubjectsByIDs(ctx context.Context, subjectIDs []bson.ObjectID) (map[bson.ObjectID]*model.Subject, error)
}

type classRepoImpl struct {
	classes  *mongo.Collection
	students *mongo.Collection
	users    *mongo.Collection
	subjects *mongo.Collection
}

// NewClassRepository creates a new ClassRepository instance.
func NewClassRepository(db *mongo.Database) ClassRepository {
	return &classRepoImpl{
		classes:  db.Collection("classes"),
		students: db.Collection("students"),
		users:    db.Collection("users"),
		subjects: db.Collection("subjects"),
	}
}

func (r *classRepoImpl) FindClassesByIDs(ctx context.Context, schoolID bson.ObjectID, classIDs []bson.ObjectID) ([]*model.Class, error) {
	if len(classIDs) == 0 {
		return []*model.Class{}, nil
	}

	filter := bson.M{
		"_id":      bson.M{"$in": classIDs},
		"schoolId": schoolID,
	}

	opts := options.Find().SetSort(bson.D{{Key: "name", Value: 1}})
	cur, err := r.classes.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var classes []*model.Class
	if err := cur.All(ctx, &classes); err != nil {
		return nil, err
	}
	return classes, nil
}

func (r *classRepoImpl) FindClassByID(ctx context.Context, schoolID, classID bson.ObjectID) (*model.Class, error) {
	filter := bson.M{
		"_id":      classID,
		"schoolId": schoolID,
	}

	var c model.Class
	err := r.classes.FindOne(ctx, filter).Decode(&c)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *classRepoImpl) CountStudentsByClassID(ctx context.Context, schoolID, classID bson.ObjectID) (int64, error) {
	filter := bson.M{
		"classId":  classID,
		"schoolId": schoolID,
	}
	return r.students.CountDocuments(ctx, filter)
}

func (r *classRepoImpl) FindStudentsByClassID(ctx context.Context, schoolID, classID bson.ObjectID, search string, page, limit int64) ([]*model.Student, int64, error) {
	filter := bson.M{
		"classId":  classID,
		"schoolId": schoolID,
	}

	if search != "" {
		safeSearch := regexp.QuoteMeta(search)
		nameRegex := bson.M{"$regex": safeSearch, "$options": "i"}

		// Search users by name to get userIDs
		userCur, err := r.users.Find(ctx, bson.M{
			"schoolId": schoolID,
			"role":     model.RoleStudent,
			"name":     nameRegex,
		}, options.Find().SetProjection(bson.M{"_id": 1}))

		var matchedUserIDs []bson.ObjectID
		if err == nil {
			defer userCur.Close(ctx)
			for userCur.Next(ctx) {
				var u struct {
					ID bson.ObjectID `bson:"_id"`
				}
				if err := userCur.Decode(&u); err == nil {
					matchedUserIDs = append(matchedUserIDs, u.ID)
				}
			}
		}

		orConditions := []bson.M{
			{"nis": nameRegex},
			{"nisn": nameRegex},
		}
		if len(matchedUserIDs) > 0 {
			orConditions = append(orConditions, bson.M{"userId": bson.M{"$in": matchedUserIDs}})
		}
		filter["$or"] = orConditions
	}

	total, err := r.students.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	skip := (page - 1) * limit
	opts := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.D{{Key: "nis", Value: 1}, {Key: "_id", Value: 1}})

	cur, err := r.students.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(ctx)

	var students []*model.Student
	if err := cur.All(ctx, &students); err != nil {
		return nil, 0, err
	}
	return students, total, nil
}

func (r *classRepoImpl) FindUsersByIDs(ctx context.Context, userIDs []bson.ObjectID) (map[bson.ObjectID]*model.User, error) {
	result := make(map[bson.ObjectID]*model.User)
	if len(userIDs) == 0 {
		return result, nil
	}

	filter := bson.M{
		"_id": bson.M{"$in": userIDs},
	}

	cur, err := r.users.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var u model.User
		if err := cur.Decode(&u); err != nil {
			return nil, err
		}
		result[u.ID] = &u
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *classRepoImpl) FindSubjectsByIDs(ctx context.Context, subjectIDs []bson.ObjectID) (map[bson.ObjectID]*model.Subject, error) {
	result := make(map[bson.ObjectID]*model.Subject)
	if len(subjectIDs) == 0 {
		return result, nil
	}

	filter := bson.M{
		"_id": bson.M{"$in": subjectIDs},
	}

	cur, err := r.subjects.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var s model.Subject
		if err := cur.Decode(&s); err != nil {
			return nil, err
		}
		result[s.ID] = &s
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
