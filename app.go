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
	"github.com/hashicorp/go-version"
	"github.com/spf13/cobra"
)

var benchmarkVersionMap = map[string]string{
	"cis-1.3.1": ">= 20.04",
	"cis-1.2":   ">= 18.09, < 20.04",
	"cis-1.1":   ">= 17.06, < 18.09",
	"cis-1.0":   ">= 1.13.0, < 17.06",
}

func app(cmd *cobra.Command, args []string) {

	version, err := ResolveCisVersion(benchmarkVersion, dockerVersion)
	if err != nil {
		util.ExitWithError(err)
	}

	path, err := getFilePath(version, "definitions.yaml")
	if err != nil {
		util.ExitWithError(err)
	}

	constraints, err := getConstraints()
	if err != nil {
		util.ExitWithError(err)
	}

	configPath, _ := getFilePath(version, "config.yaml")
	// Not checking for error because if file doesn't exist then it just nil and ignore.

	controls, err := getControls(path, configPath, constraints)
	if err != nil {
		util.ExitWithError(err)
	}

	summary := runControls(controls, checkList)
	err = outputResults(controls, summary)
	if err != nil {
		util.ExitWithError(err)
	}
}

// ResolveCisVersion returns final cis version to use.
func ResolveCisVersion(benchmarkVersion, dockerVersion string) (string, error) {
	var version string
	var err error

	// Benchmark flag is specify
	if benchmarkVersion != "" {

		// Check for not specify both --version and --benchmark
		if dockerVersion != "" {
			return "", fmt.Errorf("It is an error to specify both --version and --benchmark flags")
		}

		// Set given CIS benchmark version eg cis-1.2
		version = benchmarkVersion
	} else {
		// Auto-detect the version of Docker that's running
		if dockerVersion == "" {
			dockerVersion, err = getDockerVersion()
			if err != nil {
				return "", fmt.Errorf("Version check failed: %s\nAlternatively, you can specify the version with --version", err)
			}
		}
		// Set appropriate  CIS benchmark version according to docker version
		version, err = getDockerCisVersion(dockerVersion)
		if err != nil {
			return "", fmt.Errorf("Failed to get a valid CIS benchmark version for Docker version %s: %v", dockerVersion, err)
		}
	}
	return version, nil
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

func getControls(path, substitutionFile string, constraints []string) (*check.Controls, error) {
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
		substituMap, err := util.GetSubstitutionMap(substitutionData)
		if err != nil {
			return nil, err
		}
		s = util.MakeSubstitutions(s, "", substituMap)
	}
	controls, err := check.NewControls([]byte(s), constraints)
	if err != nil {
		return nil, err
	}

	return controls, err
}

// getDockerVersion returns the docker server engine version.
func getDockerVersion() (string, error) {
	cmd := exec.Command("docker", "version", "-f", "{{.Server.Version}}")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func getFilePath(version string, filename string) (string, error) {

	glog.V(2).Info(fmt.Sprintf("Looking for config for  %s", version))

	path := filepath.Join(cfgDir, version)
	file := filepath.Join(path, filename)

	glog.V(2).Info(fmt.Sprintf("Looking for config file: %s\n", file))

	_, err := os.Stat(file)
	if err != nil {
		return "", err
	}

	return file, nil
}

func getConstraints() (constraints []string, err error) {
	swarmStatus, err := GetDockerSwarm()
	if err != nil {
		glog.V(1).Info(fmt.Sprintf("Failed to get docker swarm status, %s", err))
	}

	constraints = append(constraints,
		fmt.Sprintf("docker-swarm=%s", swarmStatus),
	)

	glog.V(1).Info(fmt.Sprintf("The constraints are:, %s", constraints))
	return constraints, nil
}

// getDockerCisVersion select the correct CIS version in compare to running docker version
// TBD ocp-3.9 auto detection
func getDockerCisVersion(stringVersion string) (string, error) {
	dockerVersion, err := trimVersion(stringVersion)

	if err != nil {
		return "", err
	}

	for benchVersion, dockerConstraints := range benchmarkVersionMap {
		currConstraints, err := version.NewConstraint(dockerConstraints)
		if err != nil {
			return "", err
		}
		if currConstraints.Check(dockerVersion) {
			glog.V(2).Info(fmt.Sprintf("docker version %s satisfies constraints %s", dockerVersion, currConstraints))
			return benchVersion, nil
		}
	}

	tooOldVersion, err := version.NewConstraint("< 1.13.0")
	if err != nil {
		return "", err
	}

	// Vesions before 1.13.0 are not supported by CIS.
	if tooOldVersion.Check(dockerVersion) {
		return "", fmt.Errorf("docker version %s is too old", stringVersion)
	}

	return "", fmt.Errorf("no suitable CIS version has been found for docker version %s", stringVersion)
}

// TrimVersion function remove all Matadate or  Prerelease parts
// because constraints.Check() can't handle comparison with it.
func trimVersion(stringVersion string) (*version.Version, error) {
	currVersion, err := version.NewVersion(stringVersion)
	if err != nil {
		return nil, err
	}

	if currVersion.Metadata() != "" || currVersion.Prerelease() != "" {
		tempStrVersion := strings.Trim(strings.Replace(fmt.Sprint(currVersion.Segments()), " ", ".", -1), "[]")
		currVersion, err = version.NewVersion(tempStrVersion)
		if err != nil {
			return nil, err
		}
	}

	return currVersion, nil
}
