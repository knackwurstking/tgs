//go:build journal
// +build journal

package extensions

import "github.com/knackwurstking/tgs/extensions/journal"

func init() {
	Register = append(Register, journal.NewExtension(nil))
}
