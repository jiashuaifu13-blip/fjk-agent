package approval

import "context"

type Repository interface {
	Create(
		context.Context,
		*Approval,
	) error

	Get(
		context.Context,
		string,
	) (*Approval, error)

	Decide(
		context.Context,
		string,
		Status,
	) error

	Consume(
		context.Context,
		string,
		string,
	) error
}
