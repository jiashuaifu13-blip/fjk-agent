package permission

import "testing"

func TestLowRiskAllowed(t *testing.T) {
	engine := NewEngine()

	decision := engine.Check(
		Request{
			Required: LevelLow,
		},
	)

	if !decision.Allowed {
		t.Fatal(
			"LOW risk should be allowed",
		)
	}
}

func TestMediumRiskAllowed(t *testing.T) {
	engine := NewEngine()

	decision := engine.Check(
		Request{
			Required: LevelMedium,
		},
	)

	if !decision.Allowed {
		t.Fatal(
			"MEDIUM risk should be allowed",
		)
	}
}

func TestHighRiskRequiresApproval(
	t *testing.T,
) {
	engine := NewEngine()

	decision := engine.Check(
		Request{
			Required: LevelHigh,
		},
	)

	if decision.Allowed {
		t.Fatal(
			"HIGH risk should require approval",
		)
	}

	if !decision.RequiresApproval {
		t.Fatal(
			"HIGH risk approval flag missing",
		)
	}
}

func TestCriticalRequiresApproval(
	t *testing.T,
) {
	engine := NewEngine()

	decision := engine.Check(
		Request{
			Required: LevelCritical,
		},
	)

	if decision.Allowed {
		t.Fatal(
			"CRITICAL risk should require approval",
		)
	}

	if !decision.RequiresApproval {
		t.Fatal(
			"CRITICAL risk approval flag missing",
		)
	}
}
