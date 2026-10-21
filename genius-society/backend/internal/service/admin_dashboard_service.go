package service

import (
	"context"
	"time"

	"github.com/geniussociety/backend/internal/dto"
	"github.com/geniussociety/backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AdminDashboardService struct {
	db *mongo.Database
}

func NewAdminDashboardService(db *mongo.Database) *AdminDashboardService {
	return &AdminDashboardService{db: db}
}

func (s *AdminDashboardService) GetMetrics(ctx context.Context) (*dto.AdminDashboardResponse, error) {
	resp := &dto.AdminDashboardResponse{}

	// 1. Tanggal Hari Ini (WIB/Lokal)
	loc, _ := time.LoadLocation("Asia/Jakarta")
	if loc == nil {
		loc = time.Local
	}
	now := time.Now().In(loc)
	resp.CurrentDate = now.Format("02 January 2006")

	// 2. Tahun Ajaran & Semester Aktif
	var ay model.AcademicYear
	err := s.db.Collection("academic_years").FindOne(ctx, bson.M{"isActive": true}).Decode(&ay)
	if err == nil {
		resp.AcademicYear.Name = ay.Name
		resp.AcademicYear.Semester = ay.Semester
		resp.AcademicYear.IsActive = ay.IsActive
	} else {
		// Default jika belum ada di database
		resp.AcademicYear.Name = "2024/2025"
		resp.AcademicYear.Semester = 1
		resp.AcademicYear.IsActive = true
	}

	// 3. KPI Cards (Total Siswa, Total Guru, Total Kelas, Total Akun)
	usersCol := s.db.Collection("users")
	classesCol := s.db.Collection("classes")
	studentsCol := s.db.Collection("students")
	teachersCol := s.db.Collection("teachers")

	resp.KPI.TotalAccounts, _ = usersCol.CountDocuments(ctx, bson.M{})
	resp.KPI.TotalStudents, _ = studentsCol.CountDocuments(ctx, bson.M{})
	resp.KPI.TotalTeachers, _ = teachersCol.CountDocuments(ctx, bson.M{})
	resp.KPI.TotalClasses, _ = classesCol.CountDocuments(ctx, bson.M{})

	// 4. Distribusi Siswa per Kelas
	classCursor, err := classesCol.Find(ctx, bson.M{})
	if err == nil {
		defer classCursor.Close(ctx)
		var classes []model.Class
		if err := classCursor.All(ctx, &classes); err == nil {
			for _, c := range classes {
				count, _ := studentsCol.CountDocuments(ctx, bson.M{"classId": c.ID})
				resp.ClassDistribution = append(resp.ClassDistribution, dto.ClassStudentCount{
					ClassID:   c.ID.Hex(),
					ClassName: c.Name,
					Count:     count,
				})
			}
		}
	}

	// 5. Panel Pendukung: Kelengkapan Data
	// A. Siswa belum punya kelas (classId nil atau tidak ada)
	resp.DataCompleteness.StudentsWithoutClass, _ = studentsCol.CountDocuments(ctx, bson.M{
		"$or": []bson.M{
			{"classId": bson.M{"$exists": false}},
			{"classId": bson.NilObjectID},
		},
	})

	// B. Kelas tanpa wali kelas (homeroomTeacherId nil atau tidak ada)
	resp.DataCompleteness.ClassesWithoutWalas, _ = classesCol.CountDocuments(ctx, bson.M{
		"$or": []bson.M{
			{"homeroomTeacherId": bson.M{"$exists": false}},
			{"homeroomTeacherId": bson.NilObjectID},
		},
	})

	// C. Guru tanpa mata pelajaran
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$lookup", Value: bson.M{
			"from":         "teacher_subjects",
			"localField":   "_id",
			"foreignField": "teacherId",
			"as":           "subjects",
		}}},
		bson.D{{Key: "$match", Value: bson.M{
			"subjects": bson.M{"$size": 0},
		}}},
		bson.D{{Key: "$count", Value: "count"}},
	}
	tCursor, err := teachersCol.Aggregate(ctx, pipeline)
	if err == nil {
		defer tCursor.Close(ctx)
		var aggRes []struct {
			Count int64 `bson:"count"`
		}
		if err := tCursor.All(ctx, &aggRes); err == nil && len(aggRes) > 0 {
			resp.DataCompleteness.TeachersWithoutMapel = aggRes[0].Count
		} else {
			// Jika belum ada data relasi, hitung total guru
			resp.DataCompleteness.TeachersWithoutMapel, _ = teachersCol.CountDocuments(ctx, bson.M{})
		}
	}

	// D. Akun nonaktif
	resp.DataCompleteness.InactiveAccounts, _ = usersCol.CountDocuments(ctx, bson.M{"isActive": false})

	// 6. Tabel Akun / Siswa yang Baru Ditambahkan (Ambil 8 user terakhir)
	cur, err := usersCol.Find(ctx, bson.M{})
	if err == nil {
		defer cur.Close(ctx)
		var recentUsers []model.User
		if err := cur.All(ctx, &recentUsers); err == nil {
			// Ambil sampai 8 user terbaru
			limit := 8
			if len(recentUsers) < limit {
				limit = len(recentUsers)
			}
			for i := len(recentUsers) - 1; i >= len(recentUsers)-limit; i-- {
				u := recentUsers[i]
				createdStr := u.CreatedAt.Format("02 Jan 2006, 15:04")
				if u.CreatedAt.IsZero() {
					createdStr = "Baru saja"
				}
				resp.RecentUsers = append(resp.RecentUsers, dto.RecentUserItem{
					ID:        u.ID.Hex(),
					Name:      u.Name,
					Username:  u.Username,
					Email:     u.Email,
					Role:      u.Role,
					CreatedAt: createdStr,
					IsActive:  u.IsActive,
				})
			}
		}
	}

	return resp, nil
}
