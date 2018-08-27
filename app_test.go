package main

import (
	"testing"
)

func TestGetDockerVersion(t *testing.T) {
	_, err := getDockerVersion()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}

func TestGetDefinitionFilePath(t *testing.T) {
	t.Errorf("not implemented")
}

func TestGetControls(t *testing.T) {
	t.Errorf("not implemented")
}

func TestRunControls(t *testing.T) {
	t.Errorf("not implemented")
}
