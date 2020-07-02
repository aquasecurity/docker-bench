package main

import (
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

func TestGetControls(t *testing.T) {
	var err error
	path, err = getDefinitionFilePath(ver)
	if err != nil {
		t.Errorf("unexpected error: %s\n", err)
	}

	_, err = getControls(path, nil)
	if err != nil {
		t.Errorf("unexpected error: %s\n", err)
	}
}

func TestRunControls(t *testing.T) {
	var err error
	path, err = getDefinitionFilePath(ver)
	if err != nil {
		t.Errorf("unexpected error: %s\n", err)
	}
	control, err := getControls(path, nil)
	if err != nil {
		t.Errorf("unexpected error: %s\n", err)
	}

	// Run all checks
	_ = runControls(control, "")

	// Run only specified checks
	checkList := "1.2, 2.1"
	_ = runControls(control, checkList)
}
