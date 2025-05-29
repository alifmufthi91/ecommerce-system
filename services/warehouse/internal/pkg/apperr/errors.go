package apperr

type AppError struct {
	Code         int    `json:"code"`
	HumanMessage string `json:"message"`
	sys          error
	DebugError   *string `json:"debug,omitempty"`
}

func (e *AppError) Error() string {
	return e.sys.Error()
}
