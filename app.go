package main

import (
	"fmt"
	"github.com/aquasecurity/bench-common/util"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/aquasecurity/bench-common/runner"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

func app(cmd *cobra.Command, args []string) {
	var version string
	var err error

	// Get version of Docker benchmark to run
	if dockerVersion != "" {
		version = dockerVersion
	} else {
		version, err = getDockerVersion()
		if err != nil {
			util.ExitWithError(
				fmt.Errorf("Version check failed: %s\nAlternatively, you can specify the version with --version",
					err))
		}
	}

	path, err := getDefinitionFilePath(version)
	if err != nil {
		util.ExitWithError(err)
	}

	yamlCfg, err := ioutil.ReadFile(path)
	if err != nil {
		util.ExitWithError(err)
	}
	benchRunner, err := runner.New(yamlCfg).
		WithCheckList(checkList).Build()

	if err != nil {
		util.ExitWithError(err)
	}

	err = benchRunner.RunTestsWithOutput(jsonFmt, noRemediations, includeTestOutput)
	if err != nil {
		util.ExitWithError(err)
	}
}

// getDockerVersion returns the docker server engine version.
func getDockerVersion() (string, error) {
	cmd := exec.Command("docker", "version", "-f", "{{.Server.Version}}")
	out, err := cmd.Output()
	return string(out), err
}

func getDefinitionFilePath(version string) (string, error) {
	filename := "definitions.yaml"

	glog.V(2).Info(fmt.Sprintf("Looking for config for version %s", version))

	path := filepath.Join(cfgDir, version)
	file := filepath.Join(path, filename)

	glog.V(2).Info(fmt.Sprintf("Looking for config file: %s\n", file))

	_, err := os.Stat(file)
	if err != nil {
		return "", err
	}

	return file, nil
}
