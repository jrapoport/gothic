package config

import (
	"errors"

	"github.com/opentracing/opentracing-go"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/opentracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Tracer config
type Tracer struct {
	Enabled bool     `json:"enabled"`
	Address string   `json:"address"`
	Tags    []string `json:"tags"`
}

// StartTracer start the tracer
func (t *Tracer) StartTracer(name, version string) error {
	if !t.Enabled {
		return nil
	}
	if t.Address == "" {
		return errors.New("tracer address required")
	}
	opts := []tracer.StartOption{
		tracer.WithAgentAddr(t.Address),
	}
	if name != "" {
		opt := tracer.WithServiceName(name)
		opts = append(opts, opt)
	}
	if version != "" {
		opt := tracer.WithServiceVersion(version)
		opts = append(opts, opt)
	}
	tags := newKeyValueMap(t.Tags)
	for k, v := range tags {
		opts = append(opts, tracer.WithGlobalTag(k, v))
	}
	trc := opentracer.New(opts...)
	opentracing.SetGlobalTracer(trc)
	return nil
}
