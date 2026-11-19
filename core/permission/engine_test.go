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
			"HIGH risk should not be allowed without approval",
		)
	}

	if !decision.RequiresApproval {
		t.Fatal(
			"HIGH risk should require approval",
		)
	}
}

func TestCriticalApproved(t *testing.T) {
	engine := NewEngine()

	decision := engine.Check(
		Request{
			Required:     LevelCritical,
			UserApproved: true,
		},
	)

	if !decision.Allowed {
		t.Fatal(
			"approved CRITICAL request should be allowed",
		)
	}
}
