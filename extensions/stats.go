//go:build stats
// +build stats

package extensions

import "github.com/knackwurstking/tgs/extensions/stats"

func init() {
	Register = append(Register, stats.NewExtension(nil))
}
