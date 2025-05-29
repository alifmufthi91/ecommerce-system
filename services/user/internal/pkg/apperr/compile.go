package apperr

import (
	"net/http"
	"strings"
)

func CompileError(err error, lang string, debugMode bool) (int, error) {
	var (
		appError *AppError
		httpCode int
	)
	var debugErr *string
	var humanMessage []string
	if debugMode {
		errStr := err.Error()
		if len(errStr) > 0 {
			debugErr = &errStr
		}
	}

	humanMessage = strings.Split(err.Error(), "\n")

	code := ErrCode(err)
	switch code {

	case CodeValueInvalid,
		CodeHTTPBadRequest,
		CodeSQLRecordDoesNotMatch,
		CodeSQLUniqueConstraint,
		CodeSQLRecordIsExpired:
		httpCode = http.StatusBadRequest
		appError = &AppError{
			Code:         int(code),
			HumanMessage: humanMessage[0],
			sys:          err,
			DebugError:   debugErr,
		}

	case CodeHTTPNotFound:
		httpCode = http.StatusNotFound
		appError = &AppError{
			Code:         int(code),
			HumanMessage: humanMessage[0],
			sys:          err,
			DebugError:   debugErr,
		}

	case CodeHTTPUnauthorized:
		httpCode = http.StatusUnauthorized
		appError = &AppError{
			Code:         int(code),
			HumanMessage: humanMessage[0],
			sys:          err,
			DebugError:   debugErr,
		}

	case CodeHTTPUnprocessableEntity:
		httpCode = http.StatusUnprocessableEntity
		appError = &AppError{
			Code:         int(code),
			HumanMessage: humanMessage[0],
			sys:          err,
			DebugError:   debugErr,
		}

	case CodeHTTPPreconditionFailed:
		httpCode = http.StatusPreconditionFailed
		appError = &AppError{
			Code:         int(code),
			HumanMessage: humanMessage[0],
			sys:          err,
			DebugError:   debugErr,
		}

	case CodeHTTPTooManyRequests:
		httpCode = http.StatusTooManyRequests
		appError = &AppError{
			Code:         int(code),
			HumanMessage: humanMessage[0],
			sys:          err,
			DebugError:   debugErr,
		}

	case CodeHTTPForbidden:
		httpCode = http.StatusForbidden
		appError = &AppError{
			Code:         int(code),
			HumanMessage: humanMessage[0],
			sys:          err,
			DebugError:   debugErr,
		}

	default:
		httpCode = http.StatusInternalServerError
		appError = &AppError{
			Code:         int(code),
			HumanMessage: humanMessage[0],
			sys:          err,
			DebugError:   debugErr,
		}
	}

	return httpCode, appError
}
