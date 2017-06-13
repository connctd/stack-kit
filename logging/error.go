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
	User           string         `json:"user"`
	ReportLocation reportLocation `json:"reportLocation"`
	HttpRequest    httpRequest    `json:"httpRequest"`
}

func (c context) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		User           string         `json:"user"`
		ReportLocation reportLocation `json:"reportLocation"`
		HttpRequest    httpRequest    `json:"httpRequest"`
	}{c.User, c.ReportLocation, c.HttpRequest})
}

// ReportError collects all necessary information and creates the necessary key value pairs
// so that the error report can be parsed by stackdriver logging. It directly submits the error
// report to the logging subsystem and doesn't allow to add other key value pairs
func ReportError(logger log.Logger, err error, r *http.Request, subjectId string) {
	logger.Log(errorReport(err, r, subjectId)...)
}

// LogError does the same as ReportError, but returns a Logger instance and lets you log
// additional key value pairs
func LogError(logger log.Logger, err error, r *http.Request, subjectId string) log.Logger {
	return log.WithPrefix(logger, errorReport(err, r, subjectId)...)
}

func errorReport(err error, r *http.Request, subjectId string) []interface{} {
	call := stack.Caller(2)

	fileName, lineNumber := fileNameLineNumber(call)

	ctx := context{
		User: subjectId,

		ReportLocation: reportLocation{
			FilePath:     fileName,
			LineNumber:   lineNumber,
			FunctionName: fmt.Sprintf("%n", call),
		},
	}

	if r != nil {
		ctx.HttpRequest = httpRequest{
			Method:    r.Method,
			Url:       r.URL.String(),
			UserAgent: r.UserAgent(),
			Referrer:  r.Referer(),
			RemoteIp:  r.RemoteAddr,
		}
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
