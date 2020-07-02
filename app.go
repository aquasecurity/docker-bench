package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aquasecurity/bench-common/check"
	"github.com/aquasecurity/bench-common/util"
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

	path, err := getFilePath(version, "definitions.yaml")
	if err != nil {
		util.ExitWithError(err)
	}

	configPath, _ := getFilePath(version, "config.yaml")
	// Not checking for error because if file doesn't exist then it just nil and ignore.

	controls, err := getControls(path, configPath)
	if err != nil {
		util.ExitWithError(err)
	}

	summary := runControls(controls, checkList)
	err = outputResults(controls, summary)
	if err != nil {
		util.ExitWithError(err)
	}
}

func outputResults(controls *check.Controls, summary check.Summary) error {
	// if we successfully ran some tests and it's json format, ignore the warnings
	if (summary.Fail > 0 || summary.Warn > 0 || summary.Pass > 0) && jsonFmt {
		out, err := controls.JSON()
		if err != nil {
			// util.ExitWithError(fmt.Errorf("failed to output in JSON format: %v", err))
			return err
		}
		util.PrintOutput(string(out), outputFile)
	} else {
		util.PrettyPrint(controls, summary, noRemediations, includeTestOutput)
	}

	return nil
}

func runControls(controls *check.Controls, checkList string) check.Summary {
	var summary check.Summary

	if checkList != "" {
		ids := util.CleanIDs(checkList)
		summary = controls.RunChecks(ids...)
	} else {
		summary = controls.RunGroup()
	}

	return summary
}

func getControls(path string, substitutionFile string) (*check.Controls, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	s := string(data)
	if substitutionFile != "" {
		substitutionData, err := ioutil.ReadFile(substitutionFile)
		if err != nil {
			return nil, err
		}
		substituMap := util.GetSubstitutionMap(substitutionData)
		s = util.MakeSubstitutions(s, "", substituMap)
	}
	controls, err := check.NewControls([]byte(s), nil)
	if err != nil {
		return nil, err
	}

	return controls, err
}

// getDockerVersion returns the docker server engine version.
func getDockerVersion() (string, error) {
	cmd := exec.Command("docker", "version", "-f", "{{.Server.Version}}")
	out, err := cmd.Output()
	return strings.TrimSpace(string(out)), err
}

func getFilePath(version string, filename string) (string, error) {

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
