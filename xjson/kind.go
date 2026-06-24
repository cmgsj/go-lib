package xjson

//go:generate go tool stringer -trimprefix Kind -type Kind

const (
	_ Kind = iota
	KindNull
	KindBool
	KindNumber
	KindString
	KindArray
	KindObject
)

type Kind byte
