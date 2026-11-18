package teachersvc

import (
	"testing"
)

func TestCalculateOutcomeGrade(t *testing.T) {
	// Scenario 1: All items graded (2 items, scores 80 and 90) -> score 85
	scores := []float64{80.0, 90.0}
	score, prov := CalculateOutcomeGrade(scores, 2)
	if score == nil || *score != 85.0 {
		t.Fatalf("expected score 85.0, got %v", score)
	}
	if prov == nil || *prov != 85.0 {
		t.Fatalf("expected prov 85.0, got %v", prov)
	}

	// Scenario 2: Partial grading (1 of 2 graded) -> score nil, prov 80.0
	scorePart, provPart := CalculateOutcomeGrade([]float64{80.0}, 2)
	if scorePart != nil {
		t.Fatalf("expected score nil when incomplete, got %v", scorePart)
	}
	if provPart == nil || *provPart != 80.0 {
		t.Fatalf("expected prov 80.0, got %v", provPart)
	}

	// Scenario 3: 0 items -> score nil, prov nil
	scoreZero, provZero := CalculateOutcomeGrade([]float64{}, 0)
	if scoreZero != nil || provZero != nil {
		t.Fatalf("expected nil for 0 items, got %v, %v", scoreZero, provZero)
	}
}

func TestCalculateRataPP(t *testing.T) {
	sc1 := 80.0
	sc2 := 90.0

	outcomes := []OutcomeGradeCalculation{
		{Title: "PP1", ItemCount: 2, Score: &sc1},
		{Title: "PP2", ItemCount: 2, Score: &sc2},
		{Title: "PP3 Empty", ItemCount: 0, Score: nil},
	}

	rata, empty := CalculateRataPP(outcomes)
	if len(empty) != 1 || empty[0] != "PP3 Empty" {
		t.Fatalf("expected empty outcomes [PP3 Empty], got %v", empty)
	}
	if rata == nil || *rata != 85.0 {
		t.Fatalf("expected rataPP 85.0, got %v", rata)
	}

	// Incomplete outcome with items -> rataPP should be nil
	outcomes[0].Score = nil
	rataInc, _ := CalculateRataPP(outcomes)
	if rataInc != nil {
		t.Fatalf("expected nil rataPP when an outcome is incomplete, got %v", rataInc)
	}
}

func TestCalculateSemester1Score(t *testing.T) {
	rata := 85.0
	sts := 70.0

	// 70% Outcome, 30% STS: 85 * 0.7 (59.5) + 70 * 0.3 (21) = 80.5
	res := CalculateSemester1Score(&rata, &sts, 70, 30)
	if res == nil || *res != 80.5 {
		t.Fatalf("expected 80.5, got %v", res)
	}

	// 50% / 50%: 85 * 0.5 (42.5) + 70 * 0.5 (35) = 77.5
	res50 := CalculateSemester1Score(&rata, &sts, 50, 50)
	if res50 == nil || *res50 != 77.5 {
		t.Fatalf("expected 77.5, got %v", res50)
	}

	// Invalid weights (sum != 100) -> returns nil
	resInvalid := CalculateSemester1Score(&rata, &sts, 60, 30)
	if resInvalid != nil {
		t.Fatalf("expected nil for weights not summing to 100, got %v", resInvalid)
	}

	// Missing STS -> returns nil
	resNoSTS := CalculateSemester1Score(&rata, nil, 50, 50)
	if resNoSTS != nil {
		t.Fatalf("expected nil when STS missing, got %v", resNoSTS)
	}
}
