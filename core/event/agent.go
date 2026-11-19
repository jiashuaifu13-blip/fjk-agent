package event

const (
	TypeAgentRunStarted         Type = "agent.run.started"
	TypeAgentPlanCreated        Type = "agent.plan.created"
	TypeAgentStepStarted        Type = "agent.step.started"
	TypeAgentStepCompleted      Type = "agent.step.completed"
	TypeAgentRunWaitingApproval Type = "agent.run.waiting_approval"
	TypeAgentRunCompleted       Type = "agent.run.completed"
	TypeAgentRunFailed          Type = "agent.run.failed"
)
