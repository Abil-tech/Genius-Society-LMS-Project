package service

import (
	"testing"

	"github.com/geniussociety/backend/internal/model"
)

func TestIsFinalGradeClass(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{name: "12 MIPA 1", code: "XII-1", expected: true},
		{name: "Kelas XII IPS 2", code: "XII-IPS-2", expected: true},
		{name: "Kelas 12 Unggulan", code: "K12-U", expected: true},
		{name: "9A SMP", code: "IX-A", expected: true},
		{name: "Kelas 9 B", code: "IX-B", expected: true},
		{name: "Kelas 6 SD", code: "VI-A", expected: true},
		{name: "Kelas VI B", code: "VI-B", expected: true},
		{name: "10 MIPA 1", code: "X-MIPA-1", expected: false},
		{name: "Kelas XI IPS 2", code: "XI-IPS-2", expected: false},
		{name: "Kelas 7 A", code: "VII-A", expected: false},
		{name: "Kelas VII A", code: "VII-1", expected: false},
		{name: "Kelas 8 B", code: "VIII-B", expected: false},
		{name: "Kelas VIII C", code: "VIII-3", expected: false},
		{name: "Kelas 1 SD", code: "I-A", expected: false},
		{name: "11 MIPA 1", code: "XI-1", expected: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := &model.Class{
				Name: tc.name,
				Code: tc.code,
			}
			result := isFinalGradeClass(c)
			if result != tc.expected {
				t.Errorf("isFinalGradeClass(%s, %s) = %v; want %v", tc.name, tc.code, result, tc.expected)
			}
		})
	}
}
