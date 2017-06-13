package logging

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReportError(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	var buf bytes.Buffer

	SetServiceContext(ServiceContext{"Test", "0.1"})
	logger := log.NewJSONLogger(&buf)

	r, _ := http.NewRequest(http.MethodGet, "http://failing.service", nil)

	ReportError(logger, errors.New("Failure"), r, "testSubject")

	var structuredLog map[string]interface{}

	err := json.Unmarshal(buf.Bytes(), &structuredLog)
	require.Nil(err)

	serviceContext, ok := structuredLog["serviceContext"].(map[string]interface{})
	require.True(ok)
	assert.NotNil(serviceContext)
	assert.Equal("Test", serviceContext["service"])
	assert.Equal("0.1", serviceContext["version"])

	ctx, ok := structuredLog["context"].(map[string]interface{})
	require.True(ok)
	assert.Equal("testSubject", ctx["user"])

	rl, ok := ctx["reportLocation"].(map[string]interface{})
	require.True(ok)
	assert.Equal("github.com/connctd/stack-kit/logging/error_test.go", rl["filePath"])
	assert.Equal(float64(26), rl["lineNumber"])
	assert.Equal("TestReportError", rl["functionName"])

	hr, ok := ctx["httpRequest"].(map[string]interface{})
	require.True(ok)
	assert.Equal("GET", hr["method"])
	assert.Equal("http://failing.service", hr["url"])
}
