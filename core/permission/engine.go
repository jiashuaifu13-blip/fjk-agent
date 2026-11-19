package permission

type Engine struct{}

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) Check(
	request Request,
) Decision {
	switch request.Required {
	case LevelLow:
		return Decision{
			Allowed: true,
		}

	case LevelMedium:
		return Decision{
			Allowed: true,
		}

	case LevelHigh:
		if request.UserApproved {
			return Decision{
				Allowed: true,
			}
		}

		return Decision{
			Allowed:          false,
			RequiresApproval: true,
			Reason:           "high-risk tool requires approval",
		}

	case LevelCritical:
		if request.UserApproved {
			return Decision{
				Allowed: true,
			}
		}

		return Decision{
			Allowed:          false,
			RequiresApproval: true,
			Reason:           "critical tool requires approval",
		}

	default:
		return Decision{
			Allowed: false,
			Reason:  "invalid permission level",
		}
	}
}
