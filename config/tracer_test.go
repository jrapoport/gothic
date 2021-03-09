package config

import (
	"testing"

	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
)

const (
	tracerEnabled = true
	tracerAddress = "example.com:9000"
)

var testTags = map[string]string{
	"tag1": "foo",
	"tag2": "bar",
}

func TestTracer(t *testing.T) {
	runTests(t, func(t *testing.T, test testCase, c *Config) {
		tr := c.Logger.Tracer
		assert.Equal(t, tracerEnabled, tr.Enabled)
		assert.Equal(t, tracerAddress+test.mark, tr.Address)
		assert.Len(t, tr.Tags, 2)
		tags := newKeyValueMap(tr.Tags)
		for k, v := range tags {
			assert.Equal(t, testTags[k]+test.mark, v)
		}
	})
}

// tests the ENV vars are correctly taking precedence
func TestTracer_Env(t *testing.T) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			clearEnv()
			loadDotEnv(t)
			c, err := loadNormalized(test.file)
			assert.NoError(t, err)
			tr := c.Logger.Tracer
			assert.Equal(t, tracerEnabled, tr.Enabled)
			assert.Equal(t, tracerAddress, tr.Address)
			assert.Len(t, tr.Tags, 2)
			tags := newKeyValueMap(tr.Tags)
			for k, v := range tags {
				assert.Equal(t, testTags[k], v)
			}
		})
	}
}

func TestTracer_Tracer(t *testing.T) {
	runTests(t, func(t *testing.T, test testCase, c *Config) {
		tc := opentracing.GlobalTracer()
		assert.Equal(t, opentracing.NoopTracer{}, tc)
		tr := c.Logger.Tracer
		assert.Equal(t, tracerEnabled, tr.Enabled)
		assert.Equal(t, tracerAddress+test.mark, tr.Address)
		assert.Len(t, tr.Tags, 2)
		tags := newKeyValueMap(tr.Tags)
		for k, v := range tags {
			assert.Equal(t, testTags[k]+test.mark, v)
		}
	})
}

func TestTracer_StartTracer(t *testing.T) {
	t.Cleanup(func() {
		opentracing.SetGlobalTracer(opentracing.NoopTracer{})
	})
	s := serviceDefaults
	l := loggerDefaults
	tr := l.Tracer
	tr.Enabled = true
	err := tr.StartTracer(s.Name, s.Version())
	assert.Error(t, err)
	tr.Address = tracerAddress
	tr.Tags = []string{"tag1=foo", "tag2=bar"}
	err = tr.StartTracer(s.Name, s.Version())
	assert.NoError(t, err)
	has := opentracing.IsGlobalTracerRegistered()
	assert.True(t, has)
}
