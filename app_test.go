package main

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

var (
	cfgdir = "./cfg"
	ver    = "cis-1.2"
)

func TestGetDockerVersion(t *testing.T) {
	_, err := getDockerVersion()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}

// Tests all standard docker-bench defintion files
func TestGetFilePath(t *testing.T) {
	d, err := os.Open(cfgdir)
	if err != nil {
		t.Errorf("unexpected error: %s\n", err)
	}

	vers, err := d.Readdirnames(-1)
	if err != nil {
		t.Errorf("unexpected error: %s\n", err)
	}

	for _, ver := range vers {
		_, err := getFilePath(ver, "definitions.yaml")
		if err != nil {
			t.Errorf("unexpected error: %s\n", err)
		}
	}
	_, err = getFilePath("Gibrish", "definitions.yaml")
	if err.Error() != "stat cfg/Gibrish/definitions.yaml: no such file or directory" {
		t.Errorf("unexpected error: %s\n", err)
	}
	_, err = getFilePath("cis-1.2", "Gibrish.yaml")
	if err.Error() != "stat cfg/cis-1.2/Gibrish.yaml: no such file or directory" {
		t.Errorf("unexpected error: %s\n", err)
	}
}

func TestGetControls(t *testing.T) {
	tests := []struct {
		name             string
		benchmarkVersion string
		cfgFile          string
		substitutionFile string
		constraints      []string
		wantErr          bool
	}{
		{name: "Test for valid control", benchmarkVersion: "cis-1.2", cfgFile: "definitions.yaml", substitutionFile: "config.yaml", constraints: []string{"docker-swarm=inactive"}, wantErr: false},
		{name: "Test for valid control no constraints", benchmarkVersion: "cis-1.2", cfgFile: "definitions.yaml", substitutionFile: "config.yaml", wantErr: false},
		{name: "Test for invalid benchmarkVersion", benchmarkVersion: "sadf", cfgFile: "definitions.yaml", substitutionFile: "config.yaml", wantErr: true},
		{name: "Test for invalid cfg file", benchmarkVersion: "cis-1.2", cfgFile: "asd", substitutionFile: "config.yaml", wantErr: true},
		{name: "Test for invalid substitution file", benchmarkVersion: "cis-1.2", cfgFile: "definitions.yaml", substitutionFile: "asdf", constraints: []string{"docker-swarm=inactive"}, wantErr: true},
		{name: "Test for empty", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version, _ := ResolveCisVersion(tt.benchmarkVersion, "")
			path, _ := getFilePath(version, tt.cfgFile)
			substitutionFile, _ := getFilePath(version, tt.substitutionFile)
			got, err := getControls(path, substitutionFile, tt.constraints)
			if err != nil && !tt.wantErr {
				t.Errorf("unexpected error = %v\nwantErr %v\nGot: %v", err, tt.wantErr, got)
			}
		})
	}

}

func TestRunControls(t *testing.T) {
	var err error

	path, err := getFilePath(ver, "definitions.yaml")
	if err != nil {
		t.Errorf("unexpected error: %s\n", err)
	}
	configPath, err := getFilePath(ver, "config.yaml")
	if err != nil {
		t.Errorf("unexpected error: %s\n", err)
	}
	control, err := getControls(path, configPath, nil)
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
		{name: "Test for version 20.10", stringVersion: "20.10", want: "cis-1.3.1", wantErr: false},
		{name: "Test for version 20.04", stringVersion: "20.04", want: "cis-1.2", wantErr: false},
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

func Test_ResolveCisVersion(t *testing.T) {

	tests := []struct {
		name             string
		dockerVersion    string
		benchmarkVersion string
		wantErr          bool
		expect           string
	}{
		{name: "Test docker version 18.9.3", dockerVersion: "18.9.3-beta", benchmarkVersion: "", expect: "cis-1.2"},
		{name: "Test docker version 17.06.3", dockerVersion: "17.06.3", benchmarkVersion: "", expect: "cis-1.1"},
		{name: "Test docker version 15", dockerVersion: "15", benchmarkVersion: "", expect: "cis-1.0"},
		{name: "Test old docker version 1", dockerVersion: "1.11.1", benchmarkVersion: "", wantErr: true},
		{name: "Test benchmark version cis-1.2", dockerVersion: "", benchmarkVersion: "cis-1.2", expect: "cis-1.2"},
		{name: "Test benchmark version cis-1.1", dockerVersion: "", benchmarkVersion: "cis-1.1", expect: "cis-1.1"},
		{name: "Test benchmark version cis-1.0", dockerVersion: "", benchmarkVersion: "cis-1.0", expect: "cis-1.0"},
		{name: "Test both benchmark and docker version", dockerVersion: "16.3", benchmarkVersion: "cis-1.2", wantErr: true},
		// {name: "Test empty", dockerVersion: "", benchmarkVersion: "", expect: "cis-1.2"}, # TBD after set env docker version.
		{name: "Test Non exist docker version", dockerVersion: "ghdsji", benchmarkVersion: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveCisVersion(tt.benchmarkVersion, tt.dockerVersion)
			if err != nil && !tt.wantErr {
				t.Errorf("ResolveCisVersion() error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.expect) {
				t.Errorf("ResolveCisVersion() = %v, want %v", got, tt.expect)
			}
		})
	}
}

func Test_getConstraints(t *testing.T) {
	got, err := getConstraints()
	if err != nil {
		t.Errorf("unexpected error: %s\n", err)
	}

	// Currently only one constraint exist - docker-swarm
	if !strings.Contains(got[0], "docker-swarm") {
		t.Errorf("expected output to contain docker-swarm\nGot:%v\n", got)
	}

}
