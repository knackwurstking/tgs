//go:build stats
// +build stats

package extensions

import "github.com/knackwurstking/tgs/extensions/pgvis"

func init() {
	Register = append(Register, pgvis.NewExtension(nil))
}
