package permission

type Level string

const (
	LevelLow      Level = "LOW"
	LevelMedium   Level = "MEDIUM"
	LevelHigh     Level = "HIGH"
	LevelCritical Level = "CRITICAL"
)

func (l Level) Valid() bool {
	switch l {
	case
		LevelLow,
		LevelMedium,
		LevelHigh,
		LevelCritical:
		return true
	default:
		return false
	}
}
