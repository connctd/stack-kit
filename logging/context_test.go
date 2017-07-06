package logging

import (
	"bytes"
	stdcontext "context"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
)

func TestContextLogging(t *testing.T) {
	assert := assert.New(t)

	for _, ctx := range []stdcontext.Context{
		contextWithValues(map[interface{}]interface{}{"X-Request-Id": "req1", "X-Client-Id": "client1", "X-Subject-Id": "subject1"}),
	} {
		var buf bytes.Buffer

		logger := log.NewJSONLogger(&buf)

		WithContext(logger, ctx).Log("msg", "test")
		assert.Contains(buf.String(), `"X-Request-Id":"req1"`)
		assert.Contains(buf.String(), `"X-Client-Id":"client1"`)
		assert.Contains(buf.String(), `"X-Subject-Id":"subject1"`)
	}
}

func contextWithValues(vals map[interface{}]interface{}) stdcontext.Context {
	ctx := stdcontext.TODO()

	for key, val := range vals {
		ctx = stdcontext.WithValue(ctx, key, val)
	}
	return ctx
}
