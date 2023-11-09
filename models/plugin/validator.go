package plugin

import "context"

type Validator interface {
	Validate(ctx context.Context, value any) error
	ValidatePartial(ctx context.Context, value any) error
}
