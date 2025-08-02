//go:build pgpress
// +build pgpress

package extensions

import "github.com/knackwurstking/tgs/extensions/pgpress"

func init() {
	Register = append(Register, pgpress.NewExtension(nil))
}
