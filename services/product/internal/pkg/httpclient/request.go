package httpclient

import (
	"io"
	"net/url"
)

type PropRequest struct {
	URI              string
	URIParams        map[string]string
	Headers          map[string]string
	QueryParams      map[string]string
	MultiQueryParams url.Values
	AuthToken        string
	FormData         map[string]string
	MultiFormData    url.Values
	Files            map[string]string
	FileReaders      map[string]*FileReaders
	Body             interface{}
	QueryString      string
	WithRetry        bool
	MaxRetry         int
	RetryWaitTime    int
}

type FileReaders struct {
	Param    string
	FileName string
	Reader   io.Reader
}
