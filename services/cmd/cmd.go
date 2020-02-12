package cmd

import (
	"os/exec"
)

func Execute(cmd string) {
	exec.Command("bash", "-c", cmd).Output()
	// c.Start()
	// c.Wait()
}
