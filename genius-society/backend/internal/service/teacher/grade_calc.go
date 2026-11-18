package teachersvc

import (
	"math"
)

// OutcomeGradeCalculation holds calculated grades for a single Learning Outcome.
type OutcomeGradeCalculation struct {
	OutcomeID        string
	Title            string
	ItemCount        int
	GradedCount      int
	ManualCount      int
	Score            *float64 // nil if not all items are graded
	ProvisionalScore *float64 // average of graded items so far (monitoring only)
}

// StudentGradeSummary holds calculated semester grades for a student.
type StudentGradeSummary struct {
	Outcomes          []OutcomeGradeCalculation
	EmptyOutcomes     []string // outcome titles that have 0 items
	AvgOutcomes       *float64 // rataPP (nil if any outcome with items is incomplete or no outcomes)
	STS               *float64 // STS exam score
	SemesterStsScore  *float64 // final semester 1 score combining rataPP and STS
	IsComplete        bool
	MissingComponents []string // e.g. ["PP: Algoritma", "STS"]
}

// RoundHalfUp rounds a float64 to specified decimal places using half-up rule.
func RoundHalfUp(val float64, decimals int) float64 {
	pow := math.Pow(10, float64(decimals))
	if val < 0 {
		return math.Ceil(val*pow-0.5) / pow
	}
	return math.Floor(val*pow+0.5) / pow
}

// CalculateOutcomeGrade computes the score and provisional score for a single PP given item scores.
// If totalItems == 0: Score = nil, Provisional = nil.
// If gradedItems < totalItems: Score = nil (never silently 0).
func CalculateOutcomeGrade(scores []float64, totalItems int) (score *float64, provisional *float64) {
	if totalItems == 0 || len(scores) == 0 {
		return nil, nil
	}

	sum := 0.0
	for _, s := range scores {
		sum += s
	}

	prov := sum / float64(len(scores))
	provisional = &prov

	if len(scores) == totalItems {
		sc := sum / float64(totalItems)
		score = &sc
	}

	return score, provisional
}

// CalculateRataPP computes the average across all PPs that contain items.
// If any PP containing items has a nil score, rataPP is nil.
// PPs with 0 items are excluded from calculation and collected in emptyOutcomes.
func CalculateRataPP(outcomes []OutcomeGradeCalculation) (rataPP *float64, emptyOutcomes []string) {
	activeCount := 0
	sum := 0.0
	hasIncomplete := false

	for _, o := range outcomes {
		if o.ItemCount == 0 {
			emptyOutcomes = append(emptyOutcomes, o.Title)
			continue
		}
		activeCount++
		if o.Score == nil {
			hasIncomplete = true
		} else {
			sum += *o.Score
		}
	}

	if activeCount == 0 || hasIncomplete {
		return nil, emptyOutcomes
	}

	avg := sum / float64(activeCount)
	return &avg, emptyOutcomes
}

// CalculateSemester1Score combines rataPP and STS using dynamic weights from grading config.
// Weights must sum to 100. If rataPP or STS is nil, the result is nil.
// Rounding (half-up 2 decimals) is only applied to the final result.
func CalculateSemester1Score(rataPP, sts *float64, outcomeWeight, stsWeight float64) *float64 {
	if rataPP == nil || sts == nil {
		return nil
	}
	if outcomeWeight+stsWeight != 100 {
		return nil
	}

	raw := (*rataPP * (outcomeWeight / 100.0)) + (*sts * (stsWeight / 100.0))
	rounded := RoundHalfUp(raw, 2)
	return &rounded
}
