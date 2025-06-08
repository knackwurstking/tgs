///go:build ip
/// +build ip

package extensions

import "github.com/knackwurstking/tgs/extensions/ip"

func init() {
	register = append(register, ip.NewIPExtension())
}
