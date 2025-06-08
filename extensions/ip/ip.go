package ip

import "github.com/knackwurstking/tgs/pkg/extension"

type IP struct {
	extension.Extension
}

func NewIP() *IP {
	return &IP{}
}

func NewIPExtension() extension.Extension {
	return NewIP()
}

// TODO: Move stuff from `internal/botcommand/ip` to this package here
