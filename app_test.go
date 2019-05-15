package main

import (
	"github.com/aquasecurity/bench-common/runner"
	"io/ioutil"
	"os"
	"testing"
)

var (
	cfgdir = "./cfg"
	ver    = "17.06"
	path   string
)

func TestGetDockerVersion(t *testing.T) {
	_, err := getDockerVersion()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}

// Tests all standard docker-bench defintion files
func TestGetDefinitionFilePath(t *testing.T) {
	d, err := os.Open(cfgdir)
	if err != nil {
		t.Errorf("unexpected error: %s\n", err)
	}

	vers, err := d.Readdirnames(-1)
	if err != nil {
		t.Errorf("unexpected error: %s\n", err)
	}

	for _, ver := range vers {
		_, err := getDefinitionFilePath(ver)
		if err != nil {
			t.Errorf("unexpected error: %s\n", err)
		}
	}
}


func TestRunControls(t *testing.T) {
	var err error
	path, err = getDefinitionFilePath(ver)
	if err != nil {
		t.Errorf("Failed to get definition file path for version %v: %s\n",ver,  err)
	}
	var benchRunner *runner.BenchRunner

	yamlCfg, err := ioutil.ReadFile(path)
	if err != nil {
		t.Errorf("Failed to read file %v : %s\n",path, err)
	}

	benchRunner, err = runner.New(yamlCfg).Build()
	if err != nil {
		t.Errorf("Failed to create benchRunner Instance: %s\n", err)
	}
	err = benchRunner.RunTestsWithOutput(true, true, false)
	if err != nil {
		t.Errorf("Failed to run : %s\n", err)
	}

	checkList := "1.2, 2.1"
	benchRunner, err = runner.New(yamlCfg).WithCheckList(checkList).Build()
	if err != nil {
		t.Errorf("Failed to create benchRunner Instance: %s\n", err)
	}
	err = benchRunner.RunTestsWithOutput(true, true, false)
	if err != nil {
		t.Errorf("Failed to run: %s\n", err)
	}
}

