package cmd

import (
	"log"
	"os/exec"
)

func Execute(cmd string) {
	_, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		log.Println(err)
	}
}
