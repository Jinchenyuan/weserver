package transport

type NetType uint8

const (
	HTTP NetType = iota + 1
	MICRO_SERVER
	MICRO_CLIENT
)
