package rho

import (
	"context"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

var Validate = DefaultValidator

func DefaultValidator(ctx context.Context, v any) error {
	return validate.StructCtx(ctx, v)
}
