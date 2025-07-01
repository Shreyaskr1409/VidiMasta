package data

import "os/exec"

type Stream struct {
	StreamKey string
	Cmd       *exec.Cmd
}
