package validator

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	herrors "github.com/eugeneradionov/froxy/pkg/http/errors"
	"github.com/eugeneradionov/xerrors"
	v "github.com/go-playground/validator/v10"
)

var (
	validatorInstance *v.Validate
	once              sync.Once
)

// Get - initialize once and returns validator instance
func Get() *v.Validate {
	return validatorInstance
}

func FormatErrors(errs v.ValidationErrors) xerrors.XErrors {
	var validationErrs = xerrors.NewXErrsWithLen(0, len(errs))

	for i := range errs {
		validationErrs.Add(herrors.NewUnprocessableEntityError(
			fmt.Errorf("field '%s' is '%s'", errs[i].Field(), errs[i].Tag()),
			fmt.Sprintf("%s %s", errs[i].Field(), errs[i].Tag()),
		))
	}

	return validationErrs
}

func Load() (err error) { // nolint
	once.Do(func() {
		validatorInstance = v.New()

		validatorInstance.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

			if name == "-" {
				return ""
			}

			return name
		})
	})

	return err
}
