package tsrv

import "github.com/jrapoport/gothic/core"

// MockAddress is a mock host address for tests.
const MockAddress = "127.0.0.1:8080"

// MockHost is a mock host
type MockHost struct {
	HostName string
}

var _ core.Hosted = (*MockHost)(nil)

// NewMockHost returns a new mock host.
func NewMockHost(n string) *MockHost {
	return &MockHost{n}
}

// Name satisfies the Hosted interface.
func (m MockHost) Name() string {
	return m.HostName
}

// Address satisfies the Hosted interface.
func (m MockHost) Address() string {
	return MockAddress
}

// Online satisfies the Hosted interface.
func (m MockHost) Online() bool {
	return true
}

// ListenAndServe satisfies the Hosted interface.
func (m MockHost) ListenAndServe() error {
	return nil
}

// Shutdown satisfies the Hosted interface.
func (m MockHost) Shutdown() error {
	return nil
}
