package conf

import (
	"os"
	"testing"

	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	defer os.Clearenv()
	os.Exit(m.Run())
}

func TestGlobal(t *testing.T) {
	os.Setenv("GOTHIC_SITE_URL", "https://example.com")
	os.Setenv("GOTHIC_JWT_SECRET", "i-am-a-secret")

	os.Setenv("GOTHIC_DB_DRIVER", "mysql")
	os.Setenv("GOTHIC_DB_URL", "fake")
	os.Setenv("GOTHIC_REQUEST_ID", "X-Request-ID")
	gc, err := LoadConfiguration("")
	require.NoError(t, err)
	require.NotNil(t, gc)
	assert.Equal(t, "X-Request-ID", gc.RequestID)
}

func TestTracing(t *testing.T) {
	os.Setenv("GOTHIC_SITE_URL", "https://example.com")
	os.Setenv("GOTHIC_JWT_SECRET", "i-am-a-secret")

	os.Setenv("GOTHIC_DB_DRIVER", "mysql")
	os.Setenv("GOTHIC_DB_URL", "fake")
	os.Setenv("GOTHIC_TRACING_SERVICE_NAME", "identity")
	os.Setenv("GOTHIC_TRACING_PORT", "8126")
	os.Setenv("GOTHIC_TRACING_HOST", "127.0.0.1")
	os.Setenv("GOTHIC_TRACING_TAGS", "tag1:value1,tag2:value2")

	gc, _ := LoadConfiguration("")
	tc := opentracing.GlobalTracer()

	assert.Equal(t, opentracing.NoopTracer{}, tc)
	assert.Equal(t, false, gc.Tracing.Enabled)
	assert.Equal(t, "identity", gc.Tracing.ServiceName)
	assert.Equal(t, "8126", gc.Tracing.Port)
	assert.Equal(t, "127.0.0.1", gc.Tracing.Host)
	assert.Equal(t, map[string]string{"tag1": "value1", "tag2": "value2"}, gc.Tracing.Tags)
}
