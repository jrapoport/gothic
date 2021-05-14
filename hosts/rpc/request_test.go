package rpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func TestRequestContext(t *testing.T) {
	ctx := context.Background()
	rtx := RequestContext(ctx)
	assert.NotNil(t, rtx)
	assert.Equal(t, "", rtx.IPAddress())
	ctx = metadata.NewIncomingContext(ctx, metadata.New(map[string]string{
		ForwardedFor:   "127.0.0.1, 198.168.1.1",
		ReCaptchaToken: "1234",
	}))
	rtx = RequestContext(ctx)
	assert.NotNil(t, rtx)
	assert.Equal(t, "127.0.0.1", rtx.IPAddress())
	assert.Equal(t, "1234", rtx.ReCaptcha())
}
