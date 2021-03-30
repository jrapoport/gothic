package context

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/models/types/provider"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/utils"
	"github.com/stretchr/testify/assert"
)

func TestFromContext(t *testing.T) {
	t.Parallel()
	var (
		aid       = uuid.New()
		uid       = uuid.New()
		ip        = "127.0.0.1"
		prov      = provider.Google
		recaptcha = utils.SecureToken()
		sort      = store.Descending
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
	assert.Equal(t, ip, ctx.IPAddress())
	assert.Equal(t, prov, ctx.Provider())
	assert.Equal(t, recaptcha, ctx.ReCaptcha())
	assert.Equal(t, sort, ctx.Sort())
	assert.Equal(t, token, ctx.Code())
	assert.Equal(t, uid, ctx.UserID())
	assert.Equal(t, aid, ctx.AdminID())
	ctx = Background()
	assert.NotNil(t, ctx)
	ctx.SetIPAddress("")
	ctx.SetProvider("")
	ctx.SetReCaptcha("")
	ctx.SetSort(store.Ascending)
	ctx.SetCode("")
	ctx.SetUserID(uuid.Nil)
	ctx.SetAdminID(uuid.Nil)
	assert.Equal(t, "", ctx.IPAddress())
	assert.EqualValues(t, "", ctx.Provider())
	assert.Equal(t, "", ctx.ReCaptcha())
	assert.EqualValues(t, store.Ascending.String(), ctx.Sort().String())
	assert.Equal(t, "", ctx.Code())
	assert.Equal(t, uuid.Nil, ctx.UserID())
	assert.Equal(t, uuid.Nil, ctx.AdminID())
	ctx = WithValue(ctx, "foo", "bar")
	v := ctx.Value("foo")
	assert.Equal(t, "bar", v.(string))
}
