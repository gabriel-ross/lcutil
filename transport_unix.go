//go:build (linux && ignore) || darwin
// +build linux,ignore darwin

package lcutil

import (
	"os/user"
	"path"
)

var (
	DEFAULT_PEMFILE string = defaultPath()
)

func defaultPath() string {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	return path.Join(u.HomeDir, `Documents`, `riotgames.pem`)
}
