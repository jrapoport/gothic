package rpc

import (
	"errors"
	"strings"
	"testing"

	"github.com/jrapoport/gothic/test/tcore"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
)

func TestServer_RPCError(t *testing.T) {
	s, _ := tcore.Server(t, false)
	srv := NewServer(s)
	tests := []struct {
		code codes.Code
	}{
		{codes.Canceled},
		{codes.Unknown},
		{codes.InvalidArgument},
		{codes.DeadlineExceeded},
		{codes.NotFound},
		{codes.AlreadyExists},
		{codes.PermissionDenied},
		{codes.ResourceExhausted},
		{codes.FailedPrecondition},
		{codes.Aborted},
		{codes.OutOfRange},
		{codes.Unimplemented},
		{codes.Internal},
		{codes.Unavailable},
		{codes.DataLoss},
		{codes.Unauthenticated},
		{255},
	}
	for _, test := range tests {
		err := srv.RPCError(test.code, nil)
		assert.Error(t, err)
	}
	err := srv.RPCError(codes.OK, nil)
	assert.NoError(t, err)
	err = srv.RPCError(codes.Unknown, errors.New("test"))
	assert.Error(t, err)
	assert.True(t, strings.HasSuffix(err.Error(), "test"))
	assert.Equal(t, "ok", statusText(codes.OK))
}
