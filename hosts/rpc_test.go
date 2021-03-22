package hosts

/*
func configClient(t *testing.T, h core.Hosted) settings.SettingsClient {
	return tsrv.RPCClient(t, h.Address(), func(cc grpc.ClientConnInterface) interface{} {
		return settings.NewSettingsClient(cc)
	}).(settings.SettingsClient)
}

func TestRPCHost(t *testing.T) {
	t.Parallel()
	a, _, _ := tcore.API(t, false)
	// create an rcp-web host
	h := NewRPCHost(a, "127.0.0.1:0")
	require.NotNil(t, h)
	err := h.ListenAndServe()
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		return h.Online()
	},1*time.Second, 10*time.Millisecond)
	test := a.Settings()
	// unauthenticated call
	ctx := context.Background()
	ac := configClient(t, h)
	res, err := ac.Settings(ctx, &settings.SettingsRequest{})
	assert.NoError(t, err)
	assert.Equal(t, test.Status, res.Status)
	assert.Equal(t, test.Signup.Provider.Internal, res.Signup.Provider.Internal)
	// shut down
	err = h.Shutdown()
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		return !h.Online()
	},1*time.Second, 10*time.Millisecond)
}
*/
