//go:build pgvis
// +build pgvis

package extensions

import "github.com/knackwurstking/tgs/extensions/pgvis"

func init() {
	Register = append(Register, pgvis.NewExtension(nil))
}
