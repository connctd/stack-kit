package logging

import (
	"bytes"
	"errors"
	"testing"

	"github.com/go-kit/kit/log"

	"github.com/stretchr/testify/assert"
)

func TestInjectingInfoSeverity(t *testing.T) {
	logBuf := &bytes.Buffer{}
	logger := log.NewJSONLogger(logBuf)
	logger = &SeverityAwareLogger{logger}

	logger.Log("msg", "something harmless", "value", 12, "bool", true)
	assert.Contains(t, logBuf.String(), `"severity":"info"`)
}

func TestInjectingErrorLevel(t *testing.T) {
	logBuf := &bytes.Buffer{}
	logger := log.NewJSONLogger(logBuf)
	logger = &SeverityAwareLogger{logger}

	logger.Log("msg", "something harmless", "value", 12, "bool", true, "error", errors.New("Something awful"))
	assert.Contains(t, logBuf.String(), `"severity":"error"`)
}

func TestOverwritingSeverity(t *testing.T) {
	logBuf := &bytes.Buffer{}
	logger := log.NewJSONLogger(logBuf)
	logger = &SeverityAwareLogger{logger}

	logger.Log("msg", "something harmless", "value", 12, "bool", true, "error", errors.New("Something awful"), "severity", "info")
	assert.Contains(t, logBuf.String(), `"severity":"info"`)
}
