package payload

import (
	"context"
	"errors"
	"fmt"
	"strings"

	appErrors "github.com/financial_advisor/app/errors"
	"github.com/financial_advisor/app/services/transformer"
	"github.com/financial_advisor/app/services/validation"
	"github.com/go-playground/validator/v10"
)

func validate(p any) error {
	if err := validation.GetInstance().Struct(p); err != nil {
		var e validator.ValidationErrors

		switch {
		case errors.As(err, &e):
			errs := make(appErrors.SystemErrors, 0, len(e))
			for _, ee := range e {
				errs = append(errs, appErrors.NewErrorBadRequest(
					appErrors.ErrorCodeBadRequest,
					fmt.Sprintf("Invalid value for %s with tag %s: %v",
						normalizeFieldName(ee.Field()),
						strings.ToLower(ee.Tag()),
						ee.Value(),
					),
				))
			}

			return errs
		default:
			return err
		}
	}

	return nil
}

// normalizeFieldName remove suffix [0] for a field type is slice error. Ex: Users[0] -> Users
func normalizeFieldName(name string) string {
	i := strings.Index(name, "[")

	if i > 0 {
		return name[:i]
	}

	return name
}

func transform(p any) error {
	return transformer.GetInstance().Struct(context.Background(), p)
}
