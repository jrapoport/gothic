package context

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/store/types/provider"
	"github.com/jrapoport/gothic/utils"
	"github.com/stretchr/testify/assert"
)

func TestFromContext(t *testing.T) {
	var (
		aid       = uuid.New()
		uid       = uuid.New()
		ip        = "127.0.0.1"
		prov      = provider.Google
		recaptcha = utils.SecureToken()
		sort      = store.Ascending
		token     = utils.SecureToken()
	)
	ctx := Background()
	assert.NotNil(t, ctx)
	ctx.SetIPAddress(ip)
	ctx.SetProvider(prov)
	ctx.SetReCaptcha(recaptcha)
	ctx.SetSort(sort)
	ctx.SetCode(token)
	ctx.SetUserID(uid)
	ctx.SetAdminID(aid)
	assert.Equal(t, ip, ctx.GetIPAddress())
	assert.Equal(t, prov, ctx.GetProvider())
	assert.Equal(t, recaptcha, ctx.GetReCaptcha())
	assert.Equal(t, sort, ctx.GetSort())
	assert.Equal(t, token, ctx.GetCode())
	assert.Equal(t, uid, ctx.GetUserID())
	assert.Equal(t, aid, ctx.GetAdminID())
	ctx = Background()
	assert.NotNil(t, ctx)
	ctx.SetIPAddress("")
	ctx.SetProvider("")
	ctx.SetReCaptcha("")
	ctx.SetSort("")
	ctx.SetCode("")
	ctx.SetUserID(uuid.Nil)
	ctx.SetAdminID(uuid.Nil)
	assert.Equal(t, "", ctx.GetIPAddress())
	assert.EqualValues(t, "", ctx.GetProvider())
	assert.Equal(t, "", ctx.GetReCaptcha())
	assert.EqualValues(t, "", ctx.GetSort())
	assert.Equal(t, "", ctx.GetCode())
	assert.Equal(t, uuid.Nil, ctx.GetUserID())
	assert.Equal(t, uuid.Nil, ctx.GetAdminID())
	ctx = WithValue(ctx, "foo", "bar")
	v := ctx.Value("foo")
	assert.Equal(t, "bar", v.(string))
}
