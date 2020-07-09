package main

import (
	"os"
	"reflect"
	"testing"
)

var (
	cfgdir = "./cfg"
	ver    = "cis-1.2"
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

	_, err = getControls(path)
	if err != nil {
		t.Errorf("unexpected error: %s\n", err)
	}
}

func TestRunControls(t *testing.T) {
	control, err := getControls(path)
	if err != nil {
		t.Errorf("unexpected error: %s\n", err)
	}

	// Run all checks
	_ = runControls(control, "")

	// Run only specified checks
	checkList := "1.2, 2.1"
	_ = runControls(control, checkList)
}

func Test_getDockerCisVersion(t *testing.T) {
	tests := []struct {
		name          string
		stringVersion string
		want          string
		wantErr       bool
	}{
		{name: "Test for version 18.09", stringVersion: "18.09", want: "cis-1.2", wantErr: false},
		{name: "Test for version 19.3.6", stringVersion: "19.3.6", want: "cis-1.2", wantErr: false},
		{name: "Test for version 18.06", stringVersion: "18.06", want: "cis-1.1", wantErr: false},
		{name: "Test for version 18.06.0-ce", stringVersion: "18.06.0-ce", want: "cis-1.1", wantErr: false},
		{name: "Test for version 17.12", stringVersion: "17.12", want: "cis-1.1", wantErr: false},
		{name: "Test for version 17.12-beta", stringVersion: "17.12-beta", want: "cis-1.1", wantErr: false},
		{name: "Test for version 17.06", stringVersion: "17.06", want: "cis-1.1", wantErr: false},
		{name: "Test for version 17.04", stringVersion: "17.04", want: "cis-1.0", wantErr: false},
		{name: "Test for version 17.03", stringVersion: "17.03", want: "cis-1.0", wantErr: false},
		{name: "Test for version 14.04", stringVersion: "14.04", want: "cis-1.0", wantErr: false},
		{name: "Test for version 1.13.0", stringVersion: "1.13.0", want: "cis-1.0", wantErr: false},
		{name: "Test for version 1.0.1", stringVersion: "1.0.1", want: "", wantErr: true},
		{name: "Test for version asd", stringVersion: "asd", want: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getDockerCisVersion(tt.stringVersion)
			if (err != nil) != tt.wantErr {
				t.Errorf("getDockerCisVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getDockerCisVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_trimVersion(t *testing.T) {

	tests := []struct {
		name          string
		stringVersion string
		want          string
	}{
		{name: "Test remove beta", stringVersion: "18.9.3-beta", want: "18.9.3"},
		{name: "Test remove ce", stringVersion: "17.6-ce", want: "17.6.0"},
		{name: "Test remove metadata", stringVersion: "16.3+rc231", want: "16.3.0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := trimVersion(tt.stringVersion)
			if err != nil {
				t.Errorf("trimVersion() error = %v", err)
				return
			}
			if !reflect.DeepEqual(got.String(), tt.want) {
				t.Errorf("trimVersion() = %v, want %v", got.String(), tt.want)
			}
		})
	}
}
