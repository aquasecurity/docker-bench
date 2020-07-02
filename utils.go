package main

import (
	"os/exec"
	"strings"
)

func GetDockerSwarm() (platform string, err error) {
	res, err := exec.Command("sh", "-c", "docker info | grep Swarm").CombinedOutput()
	if err != nil {
		return "", err
	}

	if strings.Contains(string(res), "inactive") {
		return "inactive", nil
	}
	return "active", nil
}
