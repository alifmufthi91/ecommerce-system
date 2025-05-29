package utils

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

// Ref: https://www.convictional.com/blog/gin-validation
func ParseBindErrors(errs ...error) []string {
	var out []string
	for _, err := range errs {
		switch typedError := any(err).(type) {
		case validator.ValidationErrors:
			for _, e := range typedError {
				out = append(out, parseFieldError(e))
			}
		case *json.UnmarshalTypeError:
			out = append(out, parseMarshallingError(*typedError))
		default:
			out = append(out, err.Error())
		}
	}
	return out
}

func parseFieldError(e validator.FieldError) string {
	// workaround to the fact that the `gt|gtfield=Start` gets passed as an entire tag for some reason
	// https://github.com/go-playground/validator/issues/926
	fieldPrefix := fmt.Sprintf("field %s", e.Field())
	structField := e.StructField()
	tag := strings.Split(e.Tag(), "|")[0]
	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", fieldPrefix)

	case "required_without":
		return fmt.Sprintf("%s is required if %s is not supplied", fieldPrefix, e.Param())
	case "oneof":
		return fmt.Sprintf("%s is valid if one of %s is provided", structField, e.Param())
	case "date":
		return fmt.Sprintf("%s has invalid date format", fieldPrefix)
	case "time":
		return fmt.Sprintf("%s has invalid time format", fieldPrefix)
	case "lt", "ltfield":
		param := e.Param()
		if param == "" {
			param = time.Now().Format(time.RFC3339)
		}
		return fmt.Sprintf("%s must be less than %s", fieldPrefix, param)
	case "gt", "gtfield":
		param := e.Param()
		if param == "" {
			param = time.Now().Format(time.RFC3339)
		}
		return fmt.Sprintf("%s must be greater than %s", fieldPrefix, param)
	case "max":
		param := e.Param()
		return fmt.Sprintf("%s can't more than %s characters", fieldPrefix, param)
	default:
		return fmt.Errorf("%v", e).Error()
	}
}

func parseMarshallingError(e json.UnmarshalTypeError) string {
	return fmt.Sprintf("The field %s must be a %s", e.Field, e.Type.String())
}
