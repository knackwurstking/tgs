//go:build opmanga
// +build opmanga

package extensions

import "github.com/knackwurstking/tgs/extensions/opmanga"

func init() {
	Register = append(Register, opmanga.NewExtension(nil))
}
