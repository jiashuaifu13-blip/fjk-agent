package permission

import "github.com/fjk-dev/fjk-agent/core/tool"

func FromToolRisk(
	risk tool.RiskLevel,
) Level {
	switch risk {
	case tool.RiskLow:
		return LevelLow

	case tool.RiskMedium:
		return LevelMedium

	case tool.RiskHigh:
		return LevelHigh

	case tool.RiskCritical:
		return LevelCritical

	default:
		return ""
	}
}
