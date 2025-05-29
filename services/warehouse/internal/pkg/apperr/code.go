package apperr

import (
	"net/http"

	"github.com/palantir/stacktrace"
)

type Code = stacktrace.ErrorCode

func init() {
	stacktrace.DefaultFormat = stacktrace.FormatFull
}

var ErrCode = stacktrace.GetCode
var New = stacktrace.NewError
var NewWithCode = stacktrace.NewErrorWithCode
var RootCause = stacktrace.RootCause
var Wrap = stacktrace.Propagate
var WrapWithCode = stacktrace.PropagateWithCode
var Wrapf = stacktrace.Propagate

const (
	CodeValue       = 100
	CodeSQL         = 200
	CodeCache       = 300
	CodeHTTPClient  = 400
	Code3rdDep      = 500
	CodeHTTPHandler = 600
)

const (
	// Error On Values
	CodeValueInvalid = Code(iota + CodeValue)

	// Error On SQL
	CodeSQLBuilder = Code(iota + CodeSQL)
	CodeSQLRead
	CodeSQLRowScan
	CodeSQLCreate
	CodeSQLUpdate
	CodeSQLDelete
	CodeSQLUnlink
	CodeSQLTxBegin
	CodeSQLTxCommit
	CodeSQLPrepareStmt
	CodeSQLRecordMustExist
	CodeSQLCannotRetrieveLastInsertID
	CodeSQLCannotRetrieveAffectedRows
	CodeSQLUniqueConstraint
	CodeSQLRecordDoesNotMatch
	CodeSQLRecordIsExpired

	// Error On Cache
	CodeCacheMarshal = Code(iota + CodeCache)
	CodeCacheUnmarshal
	CodeCacheGetSimpleKey
	CodeCacheSetSimpleKey
	CodeCacheDeleteSimpleKey
	CodeCacheGetHashKey
	CodeCacheSetHashKey
	CodeCacheDeleteHashKey
	CodeCacheSetExpiration
	CodeCacheDecode
	CodeCacheLockNotAcquired
	CodeCacheLockFailed
	CodeCacheInvalidCastType

	// Error on HTTP Client
	CodeHTTPClientMarshal = Code(iota + CodeHTTPClient)
	CodeHTTPClientUnmarshal
	CodeHTTPClientErrorOnRequest
	CodeHTTPClientErrorOnReadBody

	// Error on 3rd Dep.
	CodeSMSFailure = Code(iota + Code3rdDep)
	CodeMailerFailure

	// Code HTTP Handler
	CodeHTTPBadRequest = Code(iota + CodeHTTPHandler)
	CodeHTTPNotFound
	CodeHTTPUnauthorized
	CodeHTTPInternalServerError
	CodeHTTPUnmarshal
	CodeHTTPMarshal
	CodeHTTPUnprocessableEntity
	CodeHTTPTooManyRequests
	CodeHTTPPreconditionFailed
	CodeHTTPForbidden
)

var StatusCodeToErrorCodeMap = map[int]Code{
	http.StatusBadRequest:          CodeHTTPBadRequest,
	http.StatusUnprocessableEntity: CodeHTTPUnprocessableEntity,
	http.StatusTooManyRequests:     CodeHTTPTooManyRequests,
	http.StatusPreconditionFailed:  CodeHTTPPreconditionFailed,
	http.StatusNotFound:            CodeHTTPNotFound,
	http.StatusInternalServerError: CodeHTTPInternalServerError,
	http.StatusUnauthorized:        CodeHTTPUnauthorized,
	http.StatusForbidden:           CodeHTTPForbidden,
}

func MapStatusCodeToErrorCode(code int) Code {
	errCode, ok := StatusCodeToErrorCodeMap[code]

	if !ok {
		errCode = CodeHTTPInternalServerError
	}

	return errCode
}
