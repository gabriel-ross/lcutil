//go:build (linux && ignore) || darwin
// +build linux,ignore darwin

package lcutil

import (
	"bytes"
	"os/exec"
)

// Creates a new client from an already open league of legends client using commands
// that are related to a unix based system
func NewClient() (*Client, error) {
	some_byes, err := exec.Command("ps", "-A").Output()
	if err != nil {
		return &Client{}, NotRunningErr
	}

	cmd := exec.Command("grep", "ClientUx")
	// Mimic "piping" data from a cmd
	cmd.Stdin = bytes.NewReader(some_byes)

	output, err := cmd.Output()
	if err != nil {
		return &Client{}, NotRunningErr
	}

	return newClient(output)
}
