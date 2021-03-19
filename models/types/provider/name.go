package provider

import (
	"github.com/google/uuid"
)

// Name is the name of a provider.
type Name string

// IsExternal returns true if the provider is external.
func (p Name) IsExternal() bool {
	return IsExternal(p)
}

// ID returns the name as uuid.
func (p Name) ID() uuid.UUID {
	return uuid.NewMD5(uuid.NameSpaceURL, []byte(p))
}

func (p Name) String() string {
	return string(p)
}
