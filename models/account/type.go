package account

import (
	"strings"
)

// Type is a bitmask for the provider
type Type uint32

// Type bitmask
const (
	Auth Type = 1 << iota
	Payment
	Wallet
)

// Type filters
const (
	None Type = 0x0
	All  Type = 0xffffffff
	Any       = All
)

func (t Type) Has(flag Type) bool    { return t&flag != 0 }
func (t Type) Set(flag Type) Type    { return t | flag }
func (t Type) Clear(flag Type) Type  { return t &^ flag }
func (t Type) Toggle(flag Type) Type { return t ^ flag }

func (t Type) String() string {
	if t == None {
		return ""
	}
	var flags []string
	if t.Has(Auth) {
		flags = append(flags, "auth")
	}
	if t.Has(Payment) {
		flags = append(flags, "payment")
	}
	if t.Has(Wallet) {
		flags = append(flags, "wallet")
	}
	return strings.Join(flags, ",")
}
