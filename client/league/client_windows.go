package league

import (
	"strings"

	"github.com/gabriel-ross/lcutil/client"
	"github.com/shirou/gopsutil/v3/process"
)

// Creates a new client from an already open league of legends client using commands
// that are related to a windows based system
func NewClient() (client.Client, error) {
	var invocation string
	var err error
	processes, _ := process.Processes()
	for _, process := range processes {
		exe, _ := process.Exe()
		if strings.Contains(exe, "LeagueClientUx.exe") {
			invocation, _ = process.Cmdline()
			break
		}
	}

	if err != nil {
		return &Client{}, NotRunningErr
	}

	return newClient([]byte(invocation))
}
