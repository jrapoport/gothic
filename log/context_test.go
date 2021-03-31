package log

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {
	ctx := context.Background()
	lg := FromContext(ctx)
	assert.NotNil(t, lg)
	lg = NewStdLoggerWithLevel(InfoLevel)
	ctx = WithContext(ctx, lg)
	lg = FromContext(ctx)
	assert.NotNil(t, lg)
}
