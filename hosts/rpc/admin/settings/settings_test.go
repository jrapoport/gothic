package settings

import (
	"testing"

	"github.com/golang/protobuf/jsonpb"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/segmentio/encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSettingsServer_Settings(t *testing.T) {
	t.Parallel()
	s, _ := tsrv.RPCServer(t, false)
	srv := newSettingsServer(s)
	ctx := context.Background()
	res, err := srv.Settings(ctx, nil)
	assert.NoError(t, err)
	jpb := &jsonpb.Marshaler{
		OrigName:     false,
		EnumsAsInts:  false,
		EmitDefaults: false,
		Indent:       "",
		AnyResolver:  nil,
	}
	str, err := jpb.MarshalToString(res)
	test, err := json.Marshal(s.Settings())
	require.NoError(t, err)
	assert.JSONEq(t, string(test), str)
}
