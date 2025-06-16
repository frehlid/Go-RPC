package types

import (
	"local/lib/helpers"
	"strconv"
)

// Capability to represent a "freshly" authenticated user
type UserCap uint64

func (u UserCap) String() string {
	return strconv.FormatUint(uint64(u), 10)
}

func NewUserCap(cap uint64) UserCap {
	return UserCap(cap)
}

func UnmarshalUserCap(value interface{}) UserCap {
	val, err := strconv.ParseUint(value.(string), 10, 64)
	helpers.CheckForError(err)
	return UserCap(val)
}

func MarshalUserCap(cap UserCap) interface{} {
	return strconv.FormatUint(uint64(cap), 10)
}

func NewUserCapStr(capStr string) UserCap {
	uCap, err := strconv.ParseUint(capStr, 10, 64)
	helpers.CheckForError(err)
	return UserCap(uCap)
}
