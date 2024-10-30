package admin

import (
	"testing"

	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/segmentio/encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestAdminServer_Settings(t *testing.T) {
	t.Parallel()
	s, _ := tsrv.RPCServer(t, false)
	srv := newAdminServer(s)
	_, err := srv.Settings(context.Background(), nil)
	assert.Error(t, err)
	ctx := rootContext(srv.Config())
	res, err := srv.Settings(ctx, nil)
	assert.NoError(t, err)
	jpb := &protojson.MarshalOptions{
		EmitUnpopulated: false,
		Indent:          "",
		UseProtoNames:   false,
		UseEnumNumbers:  false,
	}
	str, err := jpb.Marshal(res)
	test, err := json.Marshal(s.Settings())
	require.NoError(t, err)
	assert.JSONEq(t, string(test), string(str))
}
