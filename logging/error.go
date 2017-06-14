package logging

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-stack/stack"
)

var (
	defaultServiceContext = ServiceContext{"unset", "none"}
)

// SetServiceContext sets the service context for stackdriver error reporting. Needs to be set
// before the first calls to ReportError or LogError
func SetServiceContext(s ServiceContext) {
	defaultServiceContext = s
}

// ServiceContext defines the reported service name and version
type ServiceContext struct {
	Service string
	Version string
}

func (s ServiceContext) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Service string `json:"service"`
		Version string `json:"version"`
	}{s.Service, s.Version})
}

type httpRequest struct {
	Method             string `json:"method"`
	Url                string `json:"url"`
	UserAgent          string `json:"userAgent"`
	Referrer           string `json:"referrer"`
	ResponseStatusCode int    `json:"responseStatusCode"`
	RemoteIp           string `json:"remoteIp"`
}

type reportLocation struct {
	FilePath     string `json:"filePath"`
	LineNumber   int    `json:"lineNumber"`
	FunctionName string `json:"functionName"`
}

type context struct {
	User           string
	ReportLocation reportLocation
	HttpRequest    *httpRequest
}

func (c context) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		User           string         `json:"user,omitempty"`
		ReportLocation reportLocation `json:"reportLocation,omitempty"`
		HttpRequest    *httpRequest   `json:"httpRequest,omitempty"`
	}{c.User, c.ReportLocation, c.HttpRequest})
}

type infoFunc func(context) context

// WithSubject sets the user field of the context field in the structured log.
// This field is optional.
func WithSubject(subjectId string) infoFunc {
	return func(ctx context) context {
		ctx.User = subjectId
		return ctx
	}
}

// WithHttpRequest sets the values for httpRequest in the context field. This
// field is optional.
func WithHttpRequest(r *http.Request) infoFunc {
	return func(ctx context) context {
		if ctx.HttpRequest == nil {
			ctx.HttpRequest = &httpRequest{}
		}

		ctx.HttpRequest.Method = r.Method
		if r.URL != nil {
			ctx.HttpRequest.Url = r.URL.String()
		}
		ctx.HttpRequest.UserAgent = r.UserAgent()
		ctx.HttpRequest.Referrer = r.Referer()
		ctx.HttpRequest.RemoteIp = r.RemoteAddr

		return ctx
	}
}

// WithStatusCode sets the responseStatusCode on context.httpRequest. This field is
// optional
func WithStatusCode(status int) infoFunc {
	return func(ctx context) context {
		if ctx.HttpRequest == nil {
			ctx.HttpRequest = &httpRequest{}
		}
		ctx.HttpRequest.ResponseStatusCode = status
		return ctx
	}
}

// ReportError collects all necessary information and creates the necessary key value pairs
// so that the error report can be parsed by stackdriver logging. It directly submits the error
// report to the logging subsystem and doesn't allow to add other key value pairs
func ReportError(logger log.Logger, err error, infoFuncs ...infoFunc) {
	logger.Log(errorReport(err, infoFuncs...)...)
}

// LogError does the same as ReportError, but returns a Logger instance and lets you log
// additional key value pairs
func LogError(logger log.Logger, err error, infoFuncs ...infoFunc) log.Logger {
	return log.WithPrefix(logger, errorReport(err, infoFuncs...)...)
}

func errorReport(err error, infoFuncs ...infoFunc) []interface{} {
	call := stack.Caller(2)

	fileName, lineNumber := fileNameLineNumber(call)

	ctx := context{
		ReportLocation: reportLocation{
			FilePath:     fileName,
			LineNumber:   lineNumber,
			FunctionName: fmt.Sprintf("%n", call),
		},
	}

	for _, f := range infoFuncs {
		ctx = f(ctx)
	}

	vals := []interface{}{
		"serviceContext", defaultServiceContext,
		"context", ctx,
		"eventTime", time.Now().Format(time.RFC3339),
		"message", string(debug.Stack()),
		"error", err,
		"severity", "error",
	}
	return vals
}

func fileNameLineNumber(call stack.Call) (string, int) {
	cPath := fmt.Sprintf("%+v", call)
	parts := strings.Split(cPath, ":")
	if len(parts) != 2 {
		return "unknown", -1
	}

	lineNumber, err := strconv.Atoi(parts[1])
	if err != nil {
		lineNumber = -1
	}
	return parts[0], lineNumber
}
