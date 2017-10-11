package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/aquasecurity/bench-common/check"
	"github.com/aquasecurity/bench-common/util"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

func app(cmd *cobra.Command, args []string) {
	glog.V(2).Infof("all systems are go")

	runCheck("cfg/17.06/definitions.yaml")
}

func runCheck(f string) {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		util.ExitWithError(err)
	}

	controls, err := check.NewControls([]byte(data))
	if err != nil {
		util.ExitWithError(err)
	}

	glog.V(2).Infof("running checks in %s\n", f)

	// TODO: think about if we want to allow passing app version as parameter.
	// In that case we will have to do version verification.
	// Also think about supporting multiple versions of an app.

	summary := controls.RunGroup()

	// if we successfully ran some tests and it's json format, ignore the warnings
	if (summary.Fail > 0 || summary.Warn > 0 || summary.Pass > 0) && jsonFmt {
		out, err := controls.JSON()
		if err != nil {
			util.ExitWithError(fmt.Errorf("failed to output in JSON format: %v", err))
		}
		fmt.Println(string(out))
	} else {
		util.PrettyPrint(controls, summary)
	}

}
