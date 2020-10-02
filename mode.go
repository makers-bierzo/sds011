package sds011

type Mode byte

const (
	// The sensor read values asynchronously that you can listen setting a channel with `.OnQuery()`.
	ActiveMode Mode = 0
	// The sensor only read values when a `.Query()` is executed explicitly.
	QueryMode  Mode = 1
)
