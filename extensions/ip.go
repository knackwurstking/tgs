//go:build ip
// +build ip

package extensions

import "github.com/knackwurstking/tgs/extensions/ip"

func init() {
	Register = append(Register, ip.NewIPExtension(nil))
}
